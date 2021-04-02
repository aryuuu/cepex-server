package usecases

import "github.com/aryuuu/cepex-server/models"

type gameUsecase struct {
}

func NewGameUsecase() models.GameUsecase {
	return &gameUsecase{}
}

func (u *gameUsecase) SendMessage(connID string, message interface{}) {

}

func (u *gameUsecase) BroadcastMessage(roomID string, message interface{}) {

}
