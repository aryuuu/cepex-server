package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aryuuu/cepex-server/models"
	"github.com/aryuuu/cepex-server/models/events"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type GameRouter struct {
	Upgrader    websocket.Upgrader
	Rooms       map[string]map[*websocket.Conn]string
	GameRooms   map[string]models.Room
	GameUsecase models.GameUsecase
}

func InitGameRouter(r *mux.Router, upgrader websocket.Upgrader, guc models.GameUsecase) {
	gameRouter := &GameRouter{
		Upgrader:    upgrader,
		Rooms:       make(map[string]map[*websocket.Conn]string),
		GameRooms:   make(map[string]models.Room),
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
	log.Print(r.URL.Path)
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
			}
			return
		}
		log.Printf("gameRequest: %v", gameRequest)

		switch gameRequest.EventType {
		case "create-room":
			log.Printf("Client trying to create a new room with ID %v", roomID)

			_, ok := m.Rooms[roomID]

			if ok {
				conn.WriteJSON(events.NewCreateRoomResponse(false, roomID, models.Player{}, nil))
			} else {
				player := models.Player{
					Name:      gameRequest.ClientName,
					AvatarURL: gameRequest.AvatarURL,
					IsAlive:   true,
					PlayerID:  uuid.NewString(),
				}
				m.Rooms[roomID] = make(map[*websocket.Conn]string)
				m.Rooms[roomID][conn] = player.PlayerID

				gameRoom := models.Room{
					RoomID:      roomID,
					Capacity:    4,
					HostID:      player.PlayerID,
					IsStarted:   false,
					IsClockwise: false,
					Players:     []models.Player{},
					Deck:        models.NewDeck(),
					Count:       0,
				}

				gameRoom.AddPlayer(player)
				m.GameRooms[roomID] = gameRoom
				// m.GameRooms[roomID].Players[player.PlayerID] = player

				player.Hand = gameRoom.PickCard(2)

				log.Printf("host name: %s", player.Name)
				res := events.NewCreateRoomResponse(true, roomID, player, player.Hand)
				conn.WriteJSON(res)
			}

			break
		case "join-room":
			log.Printf("Client trying to join room %v", roomID)

			room, ok := m.Rooms[roomID]

			if ok {
				gameRoom := m.GameRooms[roomID]
				player := models.Player{
					Name:      gameRequest.ClientName,
					AvatarURL: gameRequest.AvatarURL,
					IsAlive:   true,
					PlayerID:  uuid.NewString(),
					Hand:      gameRoom.PickCard(2),
				}
				m.Rooms[roomID][conn] = player.PlayerID
				gameRoom.AddPlayer(player)

				res := events.NewJoinRoomResponse(ok, gameRoom, player.Hand)
				conn.WriteJSON(res)

				broadcast := events.NewJoinRoomBroadcast(player)
				for connection := range room {
					connection.WriteJSON(broadcast)
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
			res := events.StartGameResponse{
				EventType: "leave-room",
				StarterID: "dummyID",
			}
			conn.WriteJSON(res)
		case "play-card":
			log.Printf("Client is playing card on room %v", roomID)
			res := events.PlayCardResponse{
				EventType: "play-card",
				Success:   true,
			}
			conn.WriteJSON(res)
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
		}

		// fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(message))

		// if err = conn.WriteMessage(messageType, message); err != nil {
		// 	log.Print(err)
		// 	return
		// }
	}
}
