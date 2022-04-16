package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aryuuu/cepex-server/models"
	"github.com/gorilla/mux"
)

// ImageRouter :nodoc:
type ProfileRouter struct {
	ProfileUsecase models.ProfileUsecase
}

// InitProfileRouter :nodoc:
func InitProfileRouter(r *mux.Router, puc models.ProfileUsecase) {
	profileRouter := &ProfileRouter{
		ProfileUsecase: puc,
	}

	r.HandleFunc("/picture", profileRouter.HandleProfilePicture).Methods("POST")
	r.HandleFunc("/avatar", profileRouter.HandleAvatar).Methods("POST")
}

// HandleProfilePicture nodoc:
func (m ProfileRouter) HandleProfilePicture(w http.ResponseWriter, r *http.Request) {
	maxSize := int64(1024000)
	log.Print("POST /profile/picture")

	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		log.Print(err)
		fmt.Fprintf(w, "Image too large. Max size :%v", maxSize)
		return
	}

	file, fileHeader, err := r.FormFile("profile_picture")
	if err != nil {
		log.Print(err)
		fmt.Fprint(w, "Count not get uploaded file")
		return
	}
	defer file.Close()

	result, err := m.ProfileUsecase.UploadPicture(file, fileHeader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Failed to upload picture")
		return
	}

	body := struct {
		Data string `json:"data"`
	}{
		Data: result,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
}

func (m ProfileRouter) HandleAvatar(w http.ResponseWriter, r *http.Request) {
	maxSize := int64(1024000)
	log.Println("POST /profile/avatar")

	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		log.Print(err)
		fmt.Fprintf(w, "Image too large. Max size: %v", maxSize)
		return
	}

	avatar := r.FormValue("avatar")

	result, err := m.ProfileUsecase.UploadAvatar(avatar)
	if err != nil {
		log.Printf("ProfileRouter.HandleAvatar: error uploading avatar: %v", err)
		fmt.Fprintf(w, "Failed to upload avatar: %v", err)
		return
	}

	body := struct {
		Data string `json:"data"`
	}{
		Data: result,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(body)
}
