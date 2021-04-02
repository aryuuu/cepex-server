package main

import (
	"log"
	"net/http"

	"github.com/aryuuu/cepex-server/configs"
	"github.com/aryuuu/cepex-server/repositories"
	"github.com/aryuuu/cepex-server/routes"
	"github.com/aryuuu/cepex-server/usecases"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	// socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	r := new(mux.Router)
	r.Use(mux.CORSMethodMiddleware(r))

	// socketIOServer, err := socketio.NewServer(nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	s3Repo := repositories.NewS3Repo(configureS3())

	profileUsecase := usecases.NewProfileUsecase(s3Repo)
	gameUsecase := usecases.NewGameUsecase()

	profileRouter := r.PathPrefix("/profile").Subrouter()
	gameRouter := r.PathPrefix("/ws").Subrouter()

	routes.InitProfileRouter(profileRouter, profileUsecase)
	routes.InitGameRouter(gameRouter, upgrader, gameUsecase)

	// go socketIOServer.Serve()
	// defer socketIOServer.Close()

	// socketIOServer.OnConnect("/", func(conn socketio.Conn) error {
	// 	conn.Emit("welcome", "Bite my shiny metal ass")
	// 	log.Printf("Client %s connected to default namespace", conn.ID())
	// 	log.Printf("namespace %s", conn.Namespace())
	// 	return nil
	// })

	// r.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
	// 	allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
	// 	if origin := r.Header.Get("Origin"); origin != "" {
	// 		w.Header().Set("Access-Control-Allow-Origin", origin)
	// 		w.Header().Set("Vary", "Origin")
	// 		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
	// 		w.Header().Set("Access-Control-Allow-Credentials", "true")
	// 		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
	// 	}
	// 	if r.Method == "OPTIONS" {
	// 		return
	// 	}
	// 	log.Print("/socket.io/ ")
	// 	// log.Print(r)
	// 	log.Printf("number of connections %d", socketIOServer.Count())
	// 	socketIOServer.ServeHTTP(w, r)
	// })

	srv := &http.Server{
		Addr:    ":" + configs.Service.Port,
		Handler: r,
	}

	log.Printf("Listening on port %s...", configs.Service.Port)
	log.Fatal(srv.ListenAndServe())
}

func initService() {

}

func configureS3() *session.Session {
	s, err := session.NewSession(&aws.Config{
		Region:   aws.String("ap-south-1"),
		Endpoint: aws.String(configs.S3.ENDPOINT),
		Credentials: credentials.NewStaticCredentials(
			configs.S3.ACCESS_KEY,
			configs.S3.SECRET_KEY,
			"",
		),
	})
	if err != nil {
		log.Fatal("Failed to create s3 session")
		return nil
	}

	return s
}
