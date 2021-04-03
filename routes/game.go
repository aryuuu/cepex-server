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
			log.Printf("Client trying to create a new room with ID %v", vars["roomID"])
			roomID := vars["roomID"]

			_, ok := m.Rooms[roomID]

			if ok {
				conn.WriteJSON(events.NewCreateRoomResponse(false, roomID, models.Player{}))
			} else {
				player := models.Player{
					Name:     gameRequest.Message,
					IsAlive:  true,
					PlayerID: uuid.NewString(),
				}
				m.Rooms[roomID] = make(map[*websocket.Conn]string)
				m.Rooms[roomID][conn] = player.PlayerID

				gameRoom := models.Room{
					RoomID:      roomID,
					Capacity:    4,
					HostID:      player.PlayerID,
					IsStarted:   false,
					IsClockwise: false,
					Players:     make(map[string]models.Player),
					Deck:        make([]models.Card, 52),
					Count:       0,
				}

				m.GameRooms[roomID] = gameRoom
				m.GameRooms[roomID].Players[player.PlayerID] = player

				log.Printf("host name: %s", player.Name)
				res := events.NewCreateRoomResponse(true, roomID, player)
				conn.WriteJSON(res)
			}

			break
		case "join-room":
			log.Printf("Client trying to join room %v", vars["roomID"])
			roomID := vars["roomID"]
			room, ok := m.Rooms[roomID]
			res := events.NewJoinRoomResponse(ok, m.GameRooms[roomID])
			conn.WriteJSON(res)

			if ok {
				player := models.Player{
					Name:     gameRequest.Message,
					IsAlive:  true,
					PlayerID: uuid.NewString(),
				}
				m.Rooms[roomID][conn] = player.PlayerID
				m.GameRooms[roomID].Players[player.PlayerID] = player

				broadcast := events.NewJoinRoomBroadcast(player)
				for connection := range room {
					connection.WriteJSON(broadcast)
				}
			}
			break
		case "leave-room":
			log.Printf("Client trying to leave room %v", vars["roomID"])
			roomID := vars["roomID"]
			room, ok := m.Rooms[roomID]
			res := events.NewLeaveRoomResponse(true)
			conn.WriteJSON(res)

			if ok {
				broadcast := events.NewLeaveRoomBroadcast(room[conn])
				for connection := range room {
					connection.WriteJSON(broadcast)
				}
			}
			break
		case "start-game":
			log.Printf("Client trying to start game on room %v", vars["roomID"])
			res := events.StartGameResponse{
				EventType: "leave-room",
				StarterID: "dummyID",
			}
			conn.WriteJSON(res)
		case "play-card":
			log.Printf("Client is playing card on room %v", vars["roomID"])
			res := events.PlayCardResponse{
				EventType: "play-card",
				Success:   true,
			}
			conn.WriteJSON(res)
		case "chat":
			log.Printf("Client is sending chat on room %v", vars["roomID"])
			roomID := vars["roomID"]

			room, ok := m.Rooms[roomID]
			if ok {
				playerID := room[conn]
				playerName := m.GameRooms[roomID].Players[playerID].Name
				log.Printf("player %s send chat", playerName)
				broadcast := events.NewMessageBroadcast(gameRequest.Message, m.GameRooms[roomID].Players[playerID].Name)
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
