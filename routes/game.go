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
	GameUsecase models.GameUsecase
}

func InitGameRouter(r *mux.Router, upgrader websocket.Upgrader, guc models.GameUsecase) {
	gameRouter := &GameRouter{
		Upgrader:    upgrader,
		Rooms:       make(map[string]map[*websocket.Conn]string),
		GameUsecase: guc,
	}

	r.HandleFunc("/create", gameRouter.HandleCreateRoom)
	r.HandleFunc("/{roomID}", gameRouter.HandleGameEvent)
}

func (m GameRouter) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	ID := uuid.New()
	log.Printf("Create new room with ID: %s", ID.String())

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", ID.String())
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
		case "join-room":
			log.Printf("Client trying to join room %v", vars["roomID"])
			room, ok := m.Rooms[vars["roomID"]]
			res := events.JoinRoomResponse{
				EventType: "join-room",
				Success:   ok,
				NewRoom:   models.Room{},
			}
			conn.WriteJSON(res)

			if ok {
				m.Rooms[vars["roomID"]][conn] = "dummyID"
				broadcast := events.NewJoinRoomBroadcast(models.Player{})
				for connection := range room {
					connection.WriteJSON(broadcast)
				}
			}
			break
		case "create-room":
			log.Printf("Client trying to create a new room with ID %v", vars["roomID"])

			_, ok := m.Rooms[vars["roomID"]]

			if ok {
				conn.WriteJSON(events.NewCreateRoomResponse(false, models.Room{}))
			} else {
				m.Rooms[vars["roomID"]] = make(map[*websocket.Conn]string)
				m.Rooms[vars["roomID"]][conn] = "dummyID"

				log.Print(m.Rooms)

				res := events.CreateRoomResponse{
					EventType: "create-room",
					Success:   true,
					NewRoom:   models.Room{},
				}
				conn.WriteJSON(res)
			}

			break
		case "leave-room":
			log.Printf("Client trying to leave room %v", vars["roomID"])
			room, ok := m.Rooms[vars["roomID"]]
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

			room, ok := m.Rooms[vars["roomID"]]
			if ok {
				broadcast := events.NewMessageBroadcast(gameRequest.Message)
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
