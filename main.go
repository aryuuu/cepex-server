package main

import (
	"log"
	"net/http"

	"github.com/aryuuu/cepex-server/configs"
	"github.com/aryuuu/cepex-server/repositories"
	"github.com/aryuuu/cepex-server/routes"
	"github.com/aryuuu/cepex-server/usecases"

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

	httpClient := new(http.Client)

	// s3Repo := repositories.NewS3Repo(configureS3())
	imageRepository := repositories.NewImgurRepo(httpClient)

	profileUsecase := usecases.NewProfileUsecase(imageRepository)
	gameUsecase := usecases.NewGameUsecase()

	healthcheckRouter := r.PathPrefix("/healthcheck").Subrouter()
	profileRouter := r.PathPrefix("/profile").Subrouter()
	gameRouter := r.PathPrefix("/game").Subrouter()

	routes.InitHealthcheckRouter(healthcheckRouter)
	routes.InitProfileRouter(profileRouter, profileUsecase)
	routes.InitGameRouter(gameRouter, upgrader, gameUsecase)

	srv := &http.Server{
		Addr:    ":" + configs.Service.Port,
		Handler: r,
	}

	log.Printf("Listening on port %s...", configs.Service.Port)
	log.Fatal(srv.ListenAndServe())
}

// TODO: reuse S3
// func configureS3() *session.Session {
// 	s, err := session.NewSession(&aws.Config{
// 		Region:   aws.String("ap-south-1"),
// 		Endpoint: aws.String(configs.S3.ENDPOINT),
// 		Credentials: credentials.NewStaticCredentials(
// 			configs.S3.ACCESS_KEY,
// 			configs.S3.SECRET_KEY,
// 			"",
// 		),
// 	})
// 	if err != nil {
// 		log.Fatal("Failed to create s3 session")
// 		return nil
// 	}

// 	return s
// }
