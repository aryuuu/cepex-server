package main

import (
	"log"
	"net/http"

	"github.com/aryuuu/cepex-server/routes"
	"github.com/gorilla/mux"
)

func main() {
	r := new(mux.Router)

	routes.InitImageRouter(r)

	srv := &http.Server{
		Handler: r,
		Addr:    ":3000",
	}

	log.Print("Listening on port 3000...")
	log.Fatal(srv.ListenAndServe())
}
