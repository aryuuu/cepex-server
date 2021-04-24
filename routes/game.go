package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aryuuu/cepex-server/models"
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

	m.GameUsecase.Connect(conn, roomID)
}
