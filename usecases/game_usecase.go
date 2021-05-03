package usecases

import (
	"log"
	"math/rand"

	"github.com/aryuuu/cepex-server/models"
	"github.com/aryuuu/cepex-server/models/events"
	"github.com/aryuuu/cepex-server/utils/common"
	"github.com/gorilla/websocket"
)

type gameUsecase struct {
	Rooms     map[string]map[*websocket.Conn]string
	GameRooms map[string]*models.Room
}

func NewGameUsecase() models.GameUsecase {
	return &gameUsecase{
		Rooms:     make(map[string]map[*websocket.Conn]string),
		GameRooms: make(map[string]*models.Room),
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
				u.kickPlayer(conn, roomID)
			}
			return
		}
		log.Printf("gameRequest: %v", gameRequest)

		switch gameRequest.EventType {
		case "create-room":
			log.Printf("Client trying to create a new room with ID %v", roomID)

			_, ok := u.Rooms[roomID]

			if ok {
				conn.WriteJSON(events.NewCreateRoomResponse(false, roomID, &models.Player{}, nil))
			} else {
				player := models.NewPlayer(gameRequest.ClientName, gameRequest.AvatarURL)
				u.Rooms[roomID] = make(map[*websocket.Conn]string)
				u.Rooms[roomID][conn] = player.PlayerID

				gameRoom := models.NewRoom(roomID, player.PlayerID, 4)

				gameRoom.AddPlayer(player)
				u.GameRooms[roomID] = gameRoom
				pickedCard := gameRoom.PickCard(2)
				player.AddHand(pickedCard)

				res := events.NewCreateRoomResponse(true, roomID, player, player.Hand)
				conn.WriteJSON(res)
			}

			break
		case "join-room":
			log.Printf("Client trying to join room %v", roomID)

			room, ok := u.Rooms[roomID]

			if ok {
				gameRoom := u.GameRooms[roomID]
				player := models.NewPlayer(gameRequest.ClientName, gameRequest.AvatarURL)
				player.Hand = gameRoom.PickCard(2)

				room[conn] = player.PlayerID
				gameRoom.AddPlayer(player)

				res := events.NewJoinRoomResponse(ok, gameRoom, player.Hand)
				conn.WriteJSON(res)

				broadcast := events.NewJoinRoomBroadcast(player)
				for connection, playerID := range room {
					if playerID != player.PlayerID {
						connection.WriteJSON(broadcast)
					}
				}
			} else {
				res := events.NewJoinRoomResponse(ok, u.GameRooms[roomID], nil)
				conn.WriteJSON(res)
			}
			break
		case "leave-room":
			u.kickPlayer(conn, roomID)

			break
		case "start-game":
			log.Printf("Client trying to start game on room %v", roomID)
			room := u.Rooms[roomID]
			gameRoom := u.GameRooms[roomID]
			playerID := u.Rooms[roomID][conn]

			if playerID != gameRoom.HostID {
				res := events.NewStartGameResponse(false)
				conn.WriteJSON(res)

			} else {
				gameRoom.StartGame()
				starterIndex := rand.Intn(len(gameRoom.Players))
				gameRoom.TurnID = gameRoom.Players[starterIndex].PlayerID
				notifContent := "game started, " + gameRoom.Players[starterIndex].Name + "'s turn"

				notification := events.NewNotificationBroadcast(notifContent)
				res := events.NewStartGameBroadcast(starterIndex)
				for connection := range room {
					connection.WriteJSON(res)
					connection.WriteJSON(notification)
				}
			}
			break

		case "play-card":
			gameRoom := u.GameRooms[roomID]
			playerID := u.Rooms[roomID][conn]
			log.Printf("game turnID: %v, playerID: %v", gameRoom.TurnID, playerID)
			if !gameRoom.IsStarted {
				log.Printf("game is not started")
				res := events.NewPlayCardResponse(false, nil)
				conn.WriteJSON(res)
				break
			}

			if gameRoom.TurnID != playerID {
				log.Printf("its not your turn yet")
				res := events.NewPlayCardResponse(false, nil)
				conn.WriteJSON(res)
				break
			}

			playerIndex := -1
			for i, p := range gameRoom.Players {
				if p.PlayerID == playerID {
					playerIndex = i
					break
				}
			}

			// player := gameRoom.Players[playerIndex]
			player := gameRoom.PlayerMap[playerID]

			if !player.IsAlive {
				log.Printf("this player is dead")
				res := events.NewPlayCardResponse(false, nil)
				conn.WriteJSON(res)
				break
			}

			// log.Printf("players hand before playing: %v", player.Hand)
			log.Printf("players hand before playing: %v", gameRoom.Players[playerIndex].Hand)
			playedCard := player.Hand[gameRequest.HandIndex]
			log.Printf("Players is playing: %v", playedCard)

			var res events.PlayCardResponse

			isPlayable := gameRoom.PlayCard(playedCard, gameRequest.IsAdd)
			isAvailable := player.PlayHand(gameRequest.HandIndex)

			if isPlayable && isAvailable {
				drawnCard := gameRoom.PickCard(1)
				player.AddHand(drawnCard)
				gameRoom.PutCard([]models.Card{playedCard})

				var nextPlayerIndex int
				if gameRoom.IsClockwise {
					nextPlayerIndex = (playerIndex + 1) % len(gameRoom.Players)

				} else {
					nextPlayerIndex = common.Mod(playerIndex-1, len(gameRoom.Players))
				}

				gameRoom.TurnID = gameRoom.Players[nextPlayerIndex].PlayerID
				// log.Printf("players hand after playing: %v", player.Hand)
				log.Printf("players hand after playing: %v", (gameRoom.Players)[playerIndex].Hand)

				res = events.NewPlayCardResponse(isPlayable, player.Hand)
				conn.WriteJSON(res)

				broadcast := events.NewPlayCardBroadcast(playedCard, gameRoom.Count, gameRoom.IsClockwise, nextPlayerIndex)
				for connection := range u.Rooms[roomID] {
					connection.WriteJSON(broadcast)
				}

			} else {
				log.Printf("unplayable card")
				res = events.NewPlayCardResponse(false, nil)
				conn.WriteJSON(res)
			}
			break
		case "chat":
			log.Printf("Client is sending chat on room %v", roomID)

			room, ok := u.Rooms[roomID]
			if ok {
				playerID := room[conn]
				playerName := u.GameRooms[roomID].PlayerMap[playerID].Name

				log.Printf("player %s send chat", playerName)
				broadcast := events.NewMessageBroadcast(gameRequest.Message, playerName)
				for connection := range room {
					connection.WriteJSON(broadcast)
				}
			}
		default:
			break
		}
	}
}

func (u *gameUsecase) kickPlayer(conn *websocket.Conn, roomID string) {
	log.Printf("Client trying to leave room %v", roomID)

	room, ok := u.Rooms[roomID]
	playerID := u.Rooms[roomID][conn]
	res := events.NewLeaveRoomResponse(true)
	conn.WriteJSON(res)

	if ok {
		broadcast := events.NewLeaveRoomBroadcast(playerID)
		for connection := range room {
			connection.WriteJSON(broadcast)
		}
	}

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

func (u *gameUsecase) SendMessage(connID string, message interface{}) {

}

func (u *gameUsecase) BroadcastMessage(roomID string, message interface{}) {

}
