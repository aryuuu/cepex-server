package usecases

import (
	"log"

	"github.com/aryuuu/cepex-server/models/events"
	gameModel "github.com/aryuuu/cepex-server/models/game"
	"github.com/gorilla/websocket"
)

type gameUsecase struct {
	Rooms     map[string]map[*websocket.Conn]string
	GameRooms map[string]*gameModel.Room
}

func NewGameUsecase() gameModel.GameUsecase {
	return &gameUsecase{
		Rooms:     make(map[string]map[*websocket.Conn]string),
		GameRooms: make(map[string]*gameModel.Room),
	}
}

func (u *gameUsecase) Connect(conn *websocket.Conn, roomID string) {
	for {
		var gameRequest events.GameRequest
		err := conn.ReadJSON(&gameRequest)

		if err != nil {
			log.Print(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Print("IsUnexpectedCloseError()", err)
			} else {
				log.Printf("expected close error: %v", err)
				u.kickPlayer(conn, roomID, gameRequest)
			}
			return
		}
		log.Printf("gameRequest: %v", gameRequest)

		switch gameRequest.EventType {
		case "create-room":
			u.createRoom(conn, roomID, gameRequest)
		case "join-room":
			u.joinRoom(conn, roomID, gameRequest)
		case "leave-room":
			u.kickPlayer(conn, roomID, gameRequest)
		case "kick-player":
			u.kickPlayer(conn, roomID, gameRequest)
		case "start-game":
			u.startGame(conn, roomID)
		case "play-card":
			u.playCard(conn, roomID, gameRequest)
		case "chat":
			u.broadcastChat(conn, roomID, gameRequest)
		default:
		}
	}
}

func (u *gameUsecase) createRoom(conn *websocket.Conn, roomID string, gameRequest events.GameRequest) {
	log.Printf("Client trying to create a new room with ID %v", roomID)

	_, ok := u.Rooms[roomID]

	if ok {
		conn.WriteJSON(events.NewCreateRoomResponse(false, roomID, &gameModel.Player{}))
	} else {
		player := gameModel.NewPlayer(gameRequest.ClientName, gameRequest.AvatarURL)

		u.createConnectionRoom(roomID, conn)
		u.createGameRoom(roomID, player.PlayerID)
		u.registerPlayer(roomID, conn, player)

		// u.Rooms[roomID] = make(map[*websocket.Conn]string)
		// u.Rooms[roomID][conn] = player.PlayerID

		// gameRoom := gameModel.NewRoom(roomID, player.PlayerID, 4)
		// u.GameRooms[roomID] = gameRoom
		// gameRoom.AddPlayer(player)

		// pickedCard := gameRoom.PickCard(2)
		// player.AddHand(pickedCard)

		res := events.NewCreateRoomResponse(true, roomID, player)
		conn.WriteJSON(res)
	}
}

func (u *gameUsecase) joinRoom(conn *websocket.Conn, roomID string, gameRequest events.GameRequest) {
	log.Printf("Client trying to join room %v", roomID)

	_, ok := u.Rooms[roomID]

	if ok {
		log.Printf("found room %v", roomID)
		gameRoom := u.GameRooms[roomID]
		player := gameModel.NewPlayer(gameRequest.ClientName, gameRequest.AvatarURL)
		u.registerPlayer(roomID, conn, player)

		res := events.NewJoinRoomResponse(ok, gameRoom)
		conn.WriteJSON(res)

		broadcast := events.NewJoinRoomBroadcast(player)
		u.broadcastMessage(roomID, broadcast)
	} else {
		log.Printf("room %v does not exist", roomID)
		res := events.NewJoinRoomResponse(ok, &gameModel.Room{})
		conn.WriteJSON(res)
	}
}

func (u *gameUsecase) kickPlayer(conn *websocket.Conn, roomID string, gameRequest events.GameRequest) {
	log.Printf("Client trying to leave room %v", roomID)

	var playerID string

	if gameRequest.PlayerID == "" {
		playerID = u.Rooms[roomID][conn]
	} else {
		playerID = gameRequest.PlayerID
	}

	_, ok := u.Rooms[roomID]
	res := events.NewLeaveRoomResponse(true)
	conn.WriteJSON(res)

	if ok {
		broadcast := events.NewLeaveRoomBroadcast(playerID)
		u.broadcastMessage(roomID, broadcast)
	}

	u.unregisterPlayer(roomID, conn, playerID)
}

func (u *gameUsecase) startGame(conn *websocket.Conn, roomID string) {
	log.Printf("Client trying to start game on room %v", roomID)
	gameRoom := u.GameRooms[roomID]
	playerID := u.Rooms[roomID][conn]

	if playerID != gameRoom.HostID {
		res := events.NewStartGameResponse(false)
		conn.WriteJSON(res)

	} else {
		starterIndex := gameRoom.StartGame()

		u.dealCard(roomID)

		notifContent := "game started, " + gameRoom.Players[starterIndex].Name + "'s turn"
		notification := events.NewNotificationBroadcast(notifContent)
		res := events.NewStartGameBroadcast(starterIndex)

		u.broadcastMessage(roomID, res)
		u.broadcastMessage(roomID, notification)
	}
}

func (u *gameUsecase) playCard(conn *websocket.Conn, roomID string, gameRequest events.GameRequest) {
	gameRoom := u.GameRooms[roomID]
	playerID := u.Rooms[roomID][conn]
	// log.Printf("game turnID: %v, playerID: %v", gameRoom.TurnID, playerID)
	if !gameRoom.IsStarted {
		log.Printf("game is not started")
		res := events.NewPlayCardResponse(false, nil)
		conn.WriteJSON(res)
		return
	}

	if gameRoom.TurnID != playerID {
		log.Printf("its not your turn yet")
		res := events.NewPlayCardResponse(false, nil)
		conn.WriteJSON(res)
		return
	}

	playerIndex := gameRoom.GetPlayerIndex(playerID)

	player := gameRoom.PlayerMap[playerID]

	if !player.IsAlive {
		log.Printf("this player is dead")
		res := events.NewPlayCardResponse(false, nil)
		conn.WriteJSON(res)
		return
	}

	for _, p := range gameRoom.Players {
		log.Printf("%v's card %v", p.Name, p.Hand)
	}

	playedCard := player.Hand[gameRequest.HandIndex]
	log.Printf("%v is playing: %v", player.Name, playedCard)

	var res events.PlayCardResponse

	success := true
	if err := gameRoom.PlayCard(playerID, gameRequest.HandIndex, gameRequest.IsAdd, gameRequest.PlayerID); err != nil {
		success = false
	}

	if len(player.Hand) == 0 {
		player.IsAlive = false
		deadBroadcast := events.NewDeadPlayerBroadcast(player.PlayerID)
		u.broadcastMessage(roomID, deadBroadcast)
	}

	for _, p := range gameRoom.Players {
		log.Printf("%v's card %v", p.Name, p.Hand)
	}

	if winner := gameRoom.GetWinner(); winner != "" {
		// gameRoom.IsStarted = false
		gameRoom.EndGame()
		endBroadcast := events.NewEndGameBroadcast(winner)
		u.broadcastMessage(roomID, endBroadcast)
	}

	res = events.NewPlayCardResponse(success, player.Hand)
	conn.WriteJSON(res)

	var nextPlayerIndex int
	if gameRoom.TurnID == playerID {
		nextPlayerIndex = gameRoom.NextPlayer(playerIndex)
	} else {
		nextPlayerIndex = gameRoom.GetPlayerIndex(gameRoom.TurnID)
	}

	broadcast := events.NewPlayCardBroadcast(playedCard, gameRoom.Count, gameRoom.IsClockwise, nextPlayerIndex)
	u.broadcastMessage(roomID, broadcast)
}

func (u *gameUsecase) broadcastChat(conn *websocket.Conn, roomID string, gameRequest events.GameRequest) {
	log.Printf("Client is sending chat on room %v", roomID)

	room, ok := u.Rooms[roomID]
	if ok {
		playerID := room[conn]
		playerName := u.GameRooms[roomID].PlayerMap[playerID].Name

		log.Printf("player %s send chat", playerName)
		broadcast := events.NewMessageBroadcast(gameRequest.Message, playerName)
		u.broadcastMessage(roomID, broadcast)
	}
}

func (u *gameUsecase) createConnectionRoom(roomID string, conn *websocket.Conn) {
	u.Rooms[roomID] = make(map[*websocket.Conn]string)
}

func (u *gameUsecase) createGameRoom(roomID string, hostID string) {
	gameRoom := gameModel.NewRoom(roomID, hostID, 4)
	u.GameRooms[roomID] = gameRoom
}

func (u *gameUsecase) registerPlayer(roomID string, conn *websocket.Conn, player *gameModel.Player) {
	u.Rooms[roomID][conn] = player.PlayerID
	u.GameRooms[roomID].AddPlayer(player)
}

func (u *gameUsecase) unregisterPlayer(roomID string, conn *websocket.Conn, playerID string) {
	playerIndex := -1
	for i, p := range u.GameRooms[roomID].Players {
		if p.PlayerID == playerID {
			playerIndex = i
			break
		}
	}

	gameRoom := u.GameRooms[roomID]
	gameRoom.RemovePlayer(playerIndex)
	delete(u.Rooms[roomID], conn)

	// delete empty room
	if len(u.GameRooms[roomID].Players) == 0 {
		log.Printf("delete room %v", roomID)
		delete(u.GameRooms, roomID)
		delete(u.Rooms, roomID)
	}
}

func (u *gameUsecase) broadcastMessage(roomID string, message interface{}) {
	room := u.Rooms[roomID]
	for connection := range room {
		connection.WriteJSON(message)
	}
}

func (u *gameUsecase) dealCard(roomID string) {
	room := u.Rooms[roomID]

	for connection, playerID := range room {
		player := u.GameRooms[roomID].PlayerMap[playerID]
		message := events.NewInitialHandResponse(player.Hand)
		connection.WriteJSON(message)
	}
}

// func (u *gameUsecase) SendMessage(connID string, message interface{}) {

// }

// func (u *gameUsecase) BroadcastMessage(roomID string, message interface{}) {

// }
