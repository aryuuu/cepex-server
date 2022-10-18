package usecases

import (
	"log"

	"github.com/aryuuu/cepex-server/configs"
	"github.com/aryuuu/cepex-server/models/events"
	gameModel "github.com/aryuuu/cepex-server/models/game"
	"github.com/gorilla/websocket"
)

type connection struct {
	ID    string
	Queue chan interface{}
}

type gameUsecase struct {
	Rooms       map[string]map[*websocket.Conn]*connection
	GameRooms   map[string]*gameModel.Room
	SwitchQueue chan events.SocketEvent
}

func NewConnection(ID string) *connection {
	return &connection{
		ID:    ID,
		Queue: make(chan interface{}, 256),
	}
}

func NewGameUsecase() gameModel.GameUsecase {
	return &gameUsecase{
		Rooms:       make(map[string]map[*websocket.Conn]*connection),
		GameRooms:   make(map[string]*gameModel.Room),
		SwitchQueue: make(chan events.SocketEvent, 256),
	}
}

func (u *gameUsecase) Connect(conn *websocket.Conn, roomID string) {
	for {
		var gameRequest events.GameRequest
		err := conn.ReadJSON(&gameRequest)

		if err != nil {
			log.Print(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure) {
				log.Print("IsUnexpectedCloseError()", err)
				u.kickPlayer(conn, roomID, gameRequest)
			} else {
				log.Printf("expected close error: %v", err)
			}
			return
		}
		log.Printf("gameRequest: %v", gameRequest)

		// TODO: handle bot adding event
		switch gameRequest.EventType {
		case events.CreateRoomEvent:
			u.createRoom(conn, roomID, gameRequest)
		case events.JoinRoomEvent:
			u.joinRoom(conn, roomID, gameRequest)
		case events.LeaveRoomEvent:
			u.kickPlayer(conn, roomID, gameRequest)
		case events.KickPlayerEvent:
			u.kickPlayer(conn, roomID, gameRequest)
		case events.VoteKickPlayerEvent:
			u.voteKickPlayer(conn, roomID, gameRequest)
		case events.StartGameEvent:
			u.startGame(conn, roomID)
		case events.PlayCardEvent:
			u.playCard(conn, roomID, gameRequest)
		case events.ChatEvent:
			u.broadcastChat(conn, roomID, gameRequest)
		case events.AddBotEvent:
			u.addBotPlayer(conn, roomID, gameRequest)
		case events.KickBotEvent:
			u.kickBotPlayer(conn, roomID, gameRequest)
		default:
		}
	}
}

func (u *gameUsecase) createRoom(conn *websocket.Conn, roomID string, gameRequest events.GameRequest) {
	log.Printf("Client trying to create a new room with ID %v", roomID)

	if len(u.Rooms) >= int(configs.Constant.Capacity) {
		message := events.NewCreateRoomResponse(false, roomID, nil, "Server is full")
		u.pushUnicastMessage(roomID, conn, message)
		return
	}

	_, ok := u.Rooms[roomID]

	if ok {
		message := events.NewCreateRoomResponse(false, roomID, nil, "Room already exists")
		u.pushUnicastMessage(roomID, conn, message)
		return
	}

	player := gameModel.NewPlayer(gameRequest.ClientName, gameRequest.AvatarURL)

	u.createConnectionRoom(roomID, conn)
	u.createGameRoom(roomID, player.PlayerID)
	u.registerPlayer(roomID, conn, player)

	res := events.NewCreateRoomResponse(true, roomID, player, "")
	u.pushUnicastMessage(roomID, conn, res)
}

func (u *gameUsecase) joinRoom(conn *websocket.Conn, roomID string, gameRequest events.GameRequest) {
	log.Printf("Client trying to join room %v", roomID)

	_, ok := u.Rooms[roomID]

	if !ok {
		log.Printf("room %v does not exist", roomID)
		res := events.NewJoinRoomResponse(ok, &gameModel.Room{}, "")
		conn.WriteJSON(res)
		return
	}

	log.Printf("found room %v", roomID)
	gameRoom := u.GameRooms[roomID]
	if gameRoom.IsUsernameExist(gameRequest.ClientName) {
		log.Printf("username %s already exist", gameRequest.ClientName)
		res := events.NewJoinRoomResponse(false, &gameModel.Room{}, "username already exist")
		conn.WriteJSON(res)
		return
	}

	player := gameModel.NewPlayer(gameRequest.ClientName, gameRequest.AvatarURL)
	u.registerPlayer(roomID, conn, player)

	res := events.NewJoinRoomResponse(ok, gameRoom, "")
	u.pushUnicastMessage(roomID, conn, res)

	broadcast := events.NewJoinRoomBroadcast(player)
	u.pushBroadcastMessage(roomID, broadcast)
}

func (u *gameUsecase) addBotPlayer(conn *websocket.Conn, roomID string, gameRequest events.GameRequest) {
	log.Printf("Client trying add bot to room %v", roomID)

	_, ok := u.Rooms[roomID]

	if !ok {
		log.Printf("room %v does not exist", roomID)
		res := events.NewJoinRoomResponse(ok, &gameModel.Room{}, "")
		conn.WriteJSON(res)
		return
	}

	// TODO: create new bot player
	// TODO: register bot to room
	// TODO: broadcast new bot player to the room

	player := gameModel.NewBotPlayer()
	// u.registerPlayer(roomID, conn, player)
	u.registerBotPlayer(roomID, player)

	// TODO: change response to adding new bot player
	gameRoom := u.GameRooms[roomID]
	res := events.NewAddBotResponse(ok, gameRoom, "")
	// u.pushMessage(false, roomID, conn, res)
	u.pushUnicastMessage(roomID, conn, res)

	broadcast := events.NewAddBotBroadcast(player)
	// u.pushMessage(true, roomID, nil, broadcast)
	u.pushBroadcastMessage(roomID, broadcast)
}

func (u *gameUsecase) kickPlayer(conn *websocket.Conn, roomID string, gameRequest events.GameRequest) {
	log.Printf("Client trying to leave room %v", roomID)

	var playerID string

	if gameRequest.PlayerID == "" {
		player := u.Rooms[roomID][conn]
		if player != nil {
			playerID = u.Rooms[roomID][conn].ID
		}
	} else {
		playerID = gameRequest.PlayerID
		room := u.GameRooms[roomID]
		if room == nil {
			res := events.NewVoteKickPlayerResponse(false)
			// u.pushMessage(false, roomID, conn, res)
			u.pushUnicastMessage(roomID, conn, res)
			return
		}

		_, ok := room.PlayerMap[playerID]
		if !ok {
			res := events.NewVoteKickPlayerResponse(false)
			u.pushUnicastMessage(roomID, conn, res)
			return
		}

		res := events.NewVoteKickPlayerResponse(true)
		u.pushUnicastMessage(roomID, conn, res)

		u.GameRooms[roomID].VoteBallot[playerID] = 0
		issuerID := u.Rooms[roomID][conn].ID
		voteKickBroadcast := events.NewVoteKickPlayerBroadcast(playerID, u.GameRooms[roomID].PlayerMap[issuerID].Name)
		u.pushBroadcastMessage(roomID, voteKickBroadcast)
		return
	}

	_, ok := u.Rooms[roomID]
	res := events.NewLeaveRoomResponse(true)
	u.pushUnicastMessage(roomID, conn, res)

	if ok {
		broadcast := events.NewLeaveRoomBroadcast(playerID)
		u.pushBroadcastMessage(roomID, broadcast)
	}

	gameRoom := u.GameRooms[roomID]

	if gameRoom == nil {
		return
	}

	// appoint new host if necessary
	if gameRoom.HostID == playerID {
		newHostID := gameRoom.NextHost()
		changeHostBroadcast := events.NewChangeHostBroadcast(newHostID)
		u.pushBroadcastMessage(roomID, changeHostBroadcast)
	}

	// choose next player if necessary
	if gameRoom.IsStarted && gameRoom.TurnID == playerID {
		playerIndex := gameRoom.GetPlayerIndex(playerID)
		nextTurnId := gameRoom.NextPlayer(playerIndex)
		newPlayerBroadcast := events.NewPlayCardBroadcast(gameModel.Card{}, gameRoom.Count, gameRoom.IsClockwise, nextTurnId)
		u.pushBroadcastMessage(roomID, newPlayerBroadcast)
	}
}

// TODO: implement this function
func (u *gameUsecase) kickBotPlayer(conn *websocket.Conn, roomID string, gameRequest events.GameRequest) {
	log.Printf("Client trying to kick bot from room %v", roomID)

	var playerID string

	if gameRequest.PlayerID == "" {
		player := u.Rooms[roomID][conn]
		if player != nil {
			playerID = u.Rooms[roomID][conn].ID
		}
	} else {
		playerID = gameRequest.PlayerID
		room := u.GameRooms[roomID]
		if room == nil {
			res := events.NewVoteKickPlayerResponse(false)
			u.pushUnicastMessage(roomID, conn, res)
			return
		}

		_, ok := room.PlayerMap[playerID]
		if !ok {
			res := events.NewVoteKickPlayerResponse(false)
			u.pushUnicastMessage(roomID, conn, res)
			return
		}

		res := events.NewVoteKickPlayerResponse(true)
		u.pushUnicastMessage(roomID, conn, res)

		u.GameRooms[roomID].VoteBallot[playerID] = 0
		issuerID := u.Rooms[roomID][conn].ID
		voteKickBroadcast := events.NewVoteKickPlayerBroadcast(playerID, u.GameRooms[roomID].PlayerMap[issuerID].Name)
		u.pushBroadcastMessage(roomID, voteKickBroadcast)
		return
	}

	_, ok := u.Rooms[roomID]
	res := events.NewLeaveRoomResponse(true)
	u.pushUnicastMessage(roomID, conn, res)

	if ok {
		broadcast := events.NewLeaveRoomBroadcast(playerID)
		u.pushBroadcastMessage(roomID, broadcast)
	}

	gameRoom := u.GameRooms[roomID]

	if gameRoom == nil {
		return
	}

	// appoint new host if necessary
	if gameRoom.HostID == playerID {
		newHostID := gameRoom.NextHost()
		changeHostBroadcast := events.NewChangeHostBroadcast(newHostID)
		u.pushBroadcastMessage(roomID, changeHostBroadcast)
	}

	// choose next player if necessary
	if gameRoom.IsStarted && gameRoom.TurnID == playerID {
		playerIndex := gameRoom.GetPlayerIndex(playerID)
		nextTurnId := gameRoom.NextPlayer(playerIndex)
		newPlayerBroadcast := events.NewPlayCardBroadcast(gameModel.Card{}, gameRoom.Count, gameRoom.IsClockwise, nextTurnId)
		u.pushBroadcastMessage(roomID, newPlayerBroadcast)
	}
}

func (u *gameUsecase) voteKickPlayer(conn *websocket.Conn, roomID string, gameRequest events.GameRequest) {
	log.Printf("Client is voting on room %v", roomID)
	gameRoom := u.GameRooms[roomID]
	// playerID := u.Rooms[roomID][conn].ID

	_, ok := gameRoom.VoteBallot[gameRequest.PlayerID]
	if !ok {
		return
	}

	if gameRequest.IsAdd {
		gameRoom.VoteBallot[gameRequest.PlayerID]++
	}
	log.Printf("current tally %v", gameRoom.VoteBallot[gameRequest.PlayerID])

	if gameRequest.IsAdd && gameRoom.VoteBallot[gameRequest.PlayerID] > len(gameRoom.Players)/2 {
		log.Printf("vote kick success, removing player")
		delete(gameRoom.VoteBallot, gameRequest.PlayerID)

		var targetConn *websocket.Conn
		connRoom := u.Rooms[roomID]

		for key, val := range connRoom {
			if val.ID == gameRequest.PlayerID {
				targetConn = key
				break
			}
		}

		if targetConn == nil {
			return
		}

		evictionNotice := events.NewLeaveRoomResponse(true)
		u.pushUnicastMessage(roomID, targetConn, evictionNotice)

		broadcast := events.NewLeaveRoomBroadcast(gameRequest.PlayerID)
		u.pushBroadcastMessage(roomID, broadcast)

		// appoint new host if necessary
		if gameRoom.HostID == gameRequest.PlayerID {
			newHostID := gameRoom.NextHost()
			changeHostBroadcast := events.NewChangeHostBroadcast(newHostID)
			u.pushBroadcastMessage(roomID, changeHostBroadcast)
		}

		// choose next player if necessary
		if gameRoom.IsStarted && gameRoom.TurnID == gameRequest.PlayerID {
			playerIndex := gameRoom.GetPlayerIndex(gameRequest.PlayerID)
			nextTurnIdx := gameRoom.NextPlayer(playerIndex)
			nextPlayerBroadcast := events.NewPlayCardBroadcast(gameModel.Card{}, gameRoom.Count, gameRoom.IsClockwise, nextTurnIdx)
			u.pushBroadcastMessage(roomID, nextPlayerBroadcast)
		}
	}
}

func (u *gameUsecase) startGame(conn *websocket.Conn, roomID string) {
	log.Printf("Client trying to start game on room %v", roomID)
	gameRoom := u.GameRooms[roomID]
	playerID := u.Rooms[roomID][conn].ID

	if playerID != gameRoom.HostID {
		res := events.NewStartGameResponse(false)
		u.pushUnicastMessage(roomID, conn, res)
		return
	}

	if len(gameRoom.Players) < 2 {
		res := events.NewStartGameResponse(false)
		u.pushUnicastMessage(roomID, conn, res)
		return
	}

	starterID := gameRoom.StartGame()

	u.dealCard(roomID)

	notifContent := "game started, " + gameRoom.PlayerMap[starterID].Name + "'s turn"
	notification := events.NewNotificationBroadcast(notifContent)
	res := events.NewStartGameBroadcast(starterID)

	u.pushBroadcastMessage(roomID, res)
	u.pushBroadcastMessage(roomID, notification)
}

func (u *gameUsecase) playCard(conn *websocket.Conn, roomID string, gameRequest events.GameRequest) {
	gameRoom := u.GameRooms[roomID]
	playerID := u.Rooms[roomID][conn].ID
	if !gameRoom.IsStarted {
		log.Printf("game is not started")
		res := events.NewPlayCardResponse(false, nil, 3, "Game is not started")
		u.pushUnicastMessage(roomID, conn, res)
		return
	}

	if gameRoom.TurnID != playerID {
		log.Printf("its not your turn yet")
		res := events.NewPlayCardResponse(false, nil, 3, "Please wait for your turn")
		u.pushUnicastMessage(roomID, conn, res)
		return
	}

	playerIndex := gameRoom.GetPlayerIndex(playerID)

	player := gameRoom.PlayerMap[playerID]

	if !player.IsAlive {
		log.Printf("this player is dead")
		res := events.NewPlayCardResponse(false, nil, 3, "You are already dead")
		u.pushUnicastMessage(roomID, conn, res)
		return
	}

	playedCard := player.Hand[gameRequest.HandIndex]
	log.Printf("%v is playing: %v", player.Name, playedCard)

	var res events.PlayCardResponse

	success := true
	if err := gameRoom.PlayCard(playerID, gameRequest.HandIndex, gameRequest.IsAdd, gameRequest.PlayerID); err != nil {
		success = false

		if !gameRequest.IsDiscard {
			player.InsertHand(playedCard, gameRequest.HandIndex)
		}
	}

	if len(player.Hand) == 0 {
		player.IsAlive = false
		deadBroadcast := events.NewDeadPlayerBroadcast(player.PlayerID)
		u.pushBroadcastMessage(roomID, deadBroadcast)
	}

	if winner := gameRoom.GetWinner(); winner != nil && winner.PlayerID != "" {
		gameRoom.EndGame(winner.PlayerID)
		endBroadcast := events.NewEndGameBroadcast(winner)
		u.pushBroadcastMessage(roomID, endBroadcast)
	}

	message := ""
	status := 0
	if !success && !gameRequest.IsDiscard {
		status = 1
		res = events.NewPlayCardResponse(false, player.Hand, status, "Try discarding hand")
		res.HandIndex = gameRequest.HandIndex
		u.pushUnicastMessage(roomID, conn, res)
		return
	}

	if !success && gameRequest.IsDiscard {
		message = "Hand discarded"
	}
	res = events.NewPlayCardResponse(success, player.Hand, status, message)
	u.pushUnicastMessage(roomID, conn, res)

	var nextPlayerId string
	if gameRoom.IsStarted {
		if gameRoom.TurnID == playerID {
			nextPlayerId = gameRoom.NextPlayer(playerIndex)
		} else {
			nextPlayerId = gameRoom.TurnID
		}
	}

	if !success {
		playedCard = gameModel.Card{}
	}
	broadcast := events.NewPlayCardBroadcast(playedCard, gameRoom.Count, gameRoom.IsClockwise, nextPlayerId)
	u.pushBroadcastMessage(roomID, broadcast)
}

func (u *gameUsecase) broadcastChat(conn *websocket.Conn, roomID string, gameRequest events.GameRequest) {
	log.Printf("Client is sending chat on room %v", roomID)

	room, ok := u.Rooms[roomID]
	if ok {
		playerID := room[conn].ID
		playerName := u.GameRooms[roomID].PlayerMap[playerID].Name

		log.Printf("player %s send chat", playerName)
		broadcast := events.NewMessageBroadcast(gameRequest.Message, playerName)
		u.pushBroadcastMessage(roomID, broadcast)
	}
}

func (u *gameUsecase) createConnectionRoom(roomID string, conn *websocket.Conn) {
	u.Rooms[roomID] = make(map[*websocket.Conn]*connection)
}

func (u *gameUsecase) createGameRoom(roomID string, hostID string) {
	gameRoom := gameModel.NewRoom(roomID, hostID, 4)
	u.GameRooms[roomID] = gameRoom
}

func (u *gameUsecase) registerBotPlayer(roomID string, player *gameModel.Player) {
	// we don't need connection for bots, skip this part
	// u.Rooms[roomID][conn] = NewConnection(player.PlayerID)
	u.GameRooms[roomID].AddPlayer(player)
	// we don't need to send any websocket events to bots either, skip this part
	// go u.writePump(conn, roomID)
}

func (u *gameUsecase) registerPlayer(roomID string, conn *websocket.Conn, player *gameModel.Player) {
	u.Rooms[roomID][conn] = NewConnection(player.PlayerID)
	u.GameRooms[roomID].AddPlayer(player)
	go u.writePump(conn, roomID)
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

func (u *gameUsecase) writePump(conn *websocket.Conn, roomID string) {
	c := u.Rooms[roomID][conn]

	defer func() {
		conn.Close()
	}()

	for {
		message := <-c.Queue
		conn.WriteJSON(message)

		if _, ok := message.(events.LeaveRoomResponse); ok {
			u.unregisterPlayer(roomID, conn, c.ID)
			return
		}
	}
}

func (u *gameUsecase) dealCard(roomID string) {
	room := u.Rooms[roomID]

	for connection, playerID := range room {
		player := u.GameRooms[roomID].PlayerMap[playerID.ID]
		message := events.NewInitialHandResponse(player.Hand)
		u.pushUnicastMessage(roomID, connection, message)
	}
}

func (u *gameUsecase) RunSwitch() {
	for {
		event := <-u.SwitchQueue
		conRoom := u.Rooms[event.RoomID]
		if conRoom == nil {
			continue
		}

		if event.EventType == events.UnicastSocketEvent {
			pConn := conRoom[event.Conn]
			if pConn == nil {
				continue
			}
			pConn.Queue <- event.Message
		} else {
			for _, con := range conRoom {
				con.Queue <- event.Message
			}
		}
	}
}

func (u *gameUsecase) pushUnicastMessage(roomID string, conn *websocket.Conn, message interface{}) {
	event := events.NewUnicastEvent(roomID, conn, message)
	u.SwitchQueue <- event
}

func (u *gameUsecase) pushBroadcastMessage(roomID string, message interface{}) {
	event := events.NewBroadcastEvent(roomID, message)
	u.SwitchQueue <- event
}
