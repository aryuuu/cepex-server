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

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	r := new(mux.Router)
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
	)
	r.Use(cors)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	s3Repo := repositories.NewS3Repo(configureS3())

	profileUsecase := usecases.NewProfileUsecase(s3Repo)
	gameUsecase := usecases.NewGameUsecase()

	profileRouter := r.PathPrefix("/profile").Subrouter()
	gameRouter := r.PathPrefix("/game").Subrouter()

	routes.InitProfileRouter(profileRouter, profileUsecase)
	routes.InitGameRouter(gameRouter, upgrader, gameUsecase)

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
