package main

import (
	"log"
	"net/http"

	"github.com/aryuuu/cepex-server/configs"
	"github.com/aryuuu/cepex-server/routes"
	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
)

func main() {
	r := new(mux.Router)
	r.Use(mux.CORSMethodMiddleware(r))

	socketIOServer, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	// socketIORouter := r.PathPrefix("/socket.io").Subrouter()
	profileRouter := r.PathPrefix("/profile").Subrouter()

	routes.InitImageRouter(profileRouter)

	go socketIOServer.Serve()
	defer socketIOServer.Close()

	socketIOServer.OnConnect("default", func(conn socketio.Conn) error {
		conn.Emit("welcome", "Bite my shiny metal ass")
		log.Printf("Client %s connected to default namespace", conn.ID())
		return nil
	})

	srv := &http.Server{
		Addr:    ":" + configs.Service.Port,
		Handler: r,
	}

	log.Print("Listening on port 3000...")
	log.Fatal(srv.ListenAndServe())
}
