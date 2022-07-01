package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type HealthcheckRouter struct{}

func InitHealthcheckRouter(r *mux.Router) {
	healthcheckRouter := HealthcheckRouter{}

	r.HandleFunc("/liveness", healthcheckRouter.handleLiveness)
}

func (hcr HealthcheckRouter) handleLiveness(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Message string `json:"message"`
	}{
		Message: "OK",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
}
