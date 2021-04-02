package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aryuuu/cepex-server/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type GameRouter struct {
	Upgrader    websocket.Upgrader
	Rooms       map[string]map[*websocket.Conn]bool
	GameUsecase models.GameUsecase
}

func InitGameRouter(r *mux.Router, upgrader websocket.Upgrader, guc models.GameUsecase) {
	gameRouter := &GameRouter{
		Upgrader:    upgrader,
		Rooms:       make(map[string]map[*websocket.Conn]bool),
		GameUsecase: guc,
	}

	r.HandleFunc("/{roomID}", gameRouter.HandleGameEvent)
}

func (m GameRouter) HandleGameEvent(w http.ResponseWriter, r *http.Request) {
	log.Print(r.URL.Path)
	conn, err := m.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Print(err)
			return
		}

		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(message))

		if err = conn.WriteMessage(messageType, message); err != nil {
			log.Print(err)
			return
		}
	}
}
