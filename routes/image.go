package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// ImageRouter :nodoc:
type ImageRouter struct{}

// InitImageRouter :nodoc:
func InitImageRouter(r *mux.Router) {
	imageRouter := &ImageRouter{}

	r.HandleFunc("/picture", imageRouter.HandleProfilePicture).Methods("POST")
}

// HandleProfilePicture nodoc:
func (*ImageRouter) HandleProfilePicture(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "POST Profile picture")
}
