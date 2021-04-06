package routes

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/aryuuu/cepex-server/models"
	"github.com/aryuuu/cepex-server/models/events"
	"github.com/aryuuu/cepex-server/utils/common"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type GameRouter struct {
	Upgrader    websocket.Upgrader
	Rooms       map[string]map[*websocket.Conn]string
	GameRooms   map[string]*models.Room
	GameUsecase models.GameUsecase
}

func InitGameRouter(r *mux.Router, upgrader websocket.Upgrader, guc models.GameUsecase) {
	gameRouter := &GameRouter{
		Upgrader:    upgrader,
		Rooms:       make(map[string]map[*websocket.Conn]string),
		GameRooms:   make(map[string]*models.Room),
		GameUsecase: guc,
	}

	r.HandleFunc("/create", gameRouter.HandleCreateRoom)
	r.HandleFunc("/{roomID}", gameRouter.HandleGameEvent)
}

func (m GameRouter) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	ID := uuid.New().String()
	log.Printf("Create new room with ID: %s", ID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", ID)
}

func (m GameRouter) HandleGameEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roomID := vars["roomID"]

	conn, err := m.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}

	for {
		var gameRequest events.GameRequest
		err = conn.ReadJSON(&gameRequest)

		if err != nil {
			log.Print(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Print("IsUnexpectedCloseError()", err)
			} else {

				log.Printf("expected close error: %v", err)
			}
			return
		}
		log.Printf("gameRequest: %v", gameRequest)

		switch gameRequest.EventType {
		case "create-room":
			log.Printf("Client trying to create a new room with ID %v", roomID)

			_, ok := m.Rooms[roomID]

			if ok {
				conn.WriteJSON(events.NewCreateRoomResponse(false, roomID, &models.Player{}, nil))
			} else {
				player := &models.Player{
					Name:      gameRequest.ClientName,
					AvatarURL: gameRequest.AvatarURL,
					IsAlive:   true,
					PlayerID:  uuid.NewString(),
					Hand:      []models.Card{},
				}
				m.Rooms[roomID] = make(map[*websocket.Conn]string)
				m.Rooms[roomID][conn] = player.PlayerID

				gameRoom := &models.Room{
					RoomID:      roomID,
					Capacity:    4,
					HostID:      player.PlayerID,
					IsStarted:   false,
					IsClockwise: false,
					Players:     []*models.Player{},
					Deck:        models.NewDeck(),
					Count:       0,
				}

				gameRoom.AddPlayer(player)
				m.GameRooms[roomID] = gameRoom
				pickedCard := gameRoom.PickCard(2)
				player.AddHand(pickedCard)

				res := events.NewCreateRoomResponse(true, roomID, player, player.Hand)
				conn.WriteJSON(res)
			}

			break
		case "join-room":
			log.Printf("Client trying to join room %v", roomID)

			room, ok := m.Rooms[roomID]

			if ok {
				gameRoom := m.GameRooms[roomID]
				player := &models.Player{
					Name:      gameRequest.ClientName,
					AvatarURL: gameRequest.AvatarURL,
					IsAlive:   true,
					PlayerID:  uuid.NewString(),
					Hand:      gameRoom.PickCard(2),
				}
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
				res := events.NewJoinRoomResponse(ok, m.GameRooms[roomID], nil)
				conn.WriteJSON(res)
			}
			break
		case "leave-room":
			log.Printf("Client trying to leave room %v", roomID)

			room, ok := m.Rooms[roomID]
			playerID := m.Rooms[roomID][conn]
			res := events.NewLeaveRoomResponse(true)
			conn.WriteJSON(res)

			if ok {
				broadcast := events.NewLeaveRoomBroadcast(playerID)
				for connection := range room {
					connection.WriteJSON(broadcast)
				}
			}

			playerIndex := -1
			for i, p := range m.GameRooms[roomID].Players {
				if p.PlayerID == playerID {
					playerIndex = i
					break
				}
			}

			gameRoom := m.GameRooms[roomID]
			gameRoom.RemovePlayer(playerIndex)
			delete(m.Rooms[roomID], conn)

			break
		case "start-game":
			log.Printf("Client trying to start game on room %v", roomID)
			room := m.Rooms[roomID]
			gameRoom := m.GameRooms[roomID]
			playerID := m.Rooms[roomID][conn]

			if playerID != gameRoom.HostID {
				res := events.NewStartGameResponse(false)
				conn.WriteJSON(res)

			} else {
				gameRoom.IsStarted = true
				starterIndex := rand.Intn(len(gameRoom.Players))
				gameRoom.TurnID = gameRoom.Players[starterIndex].PlayerID
				notifContent := "game started, " + gameRoom.Players[starterIndex].Name + "'s turn"

				notification := events.NewMessageBroadcast(notifContent, "system")
				res := events.NewStartGameBroadcast(starterIndex)
				for connection := range room {
					connection.WriteJSON(res)
					connection.WriteJSON(notification)
				}
			}
			break

		case "play-card":
			log.Printf("Client is playing card on room %v", roomID)
			gameRoom := m.GameRooms[roomID]
			playerID := m.Rooms[roomID][conn]
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

			player := gameRoom.Players[playerIndex]
			// log.Printf("players hand before playing: %v", player.Hand)
			// log.Printf("players hand before playing: %v", gameRoom.Players[playerIndex].Hand)
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
					// nextPlayerIndex = (playerIndex - 1) % len(gameRoom.Players)
					nextPlayerIndex = common.Mod(playerIndex-1, len(gameRoom.Players))
				}

				gameRoom.TurnID = gameRoom.Players[nextPlayerIndex].PlayerID
				log.Printf("players hand after playing: %v", player.Hand)
				log.Printf("players hand after playing: %v", (gameRoom.Players)[playerIndex].Hand)

				res = events.NewPlayCardResponse(isPlayable, player.Hand)
				conn.WriteJSON(res)

				broadcast := events.NewPlayCardBroadcast(playedCard, gameRoom.Count, gameRoom.IsClockwise, nextPlayerIndex)
				for connection := range m.Rooms[roomID] {
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

			room, ok := m.Rooms[roomID]
			if ok {
				playerID := room[conn]
				var playerName string
				for _, p := range m.GameRooms[roomID].Players {
					if p.PlayerID == playerID {
						playerName = p.Name
						break
					}
				}

				log.Printf("player %s send chat", playerName)
				broadcast := events.NewMessageBroadcast(gameRequest.Message, playerName)
				for connection := range room {
					connection.WriteJSON(broadcast)
				}
			}
		default:
			break
		}

		// fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(message))

		// if err = conn.WriteMessage(messageType, message); err != nil {
		// 	log.Print(err)
		// 	return
		// }
	}
}
