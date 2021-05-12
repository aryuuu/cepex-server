package game

import (
	"errors"

	"github.com/google/uuid"
)

// Player :nodoc:
type Player struct {
	PlayerID  string `json:"id_player,omitempty"`
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatar_url"`
	IsAlive   bool   `json:"is_alive"`
	Hand      []Card `json:"-"`
}

func NewPlayer(name, avatarUrl string) *Player {
	return &Player{
		Name:      name,
		AvatarURL: avatarUrl,
		PlayerID:  uuid.NewString(),
		IsAlive:   false,
		Hand:      []Card{},
	}
}

func (p *Player) PlayHand(index int) error {
	// log.Println(p.Name, "'s hand ", p.Hand)
	if index >= len(p.Hand) {
		return errors.New("Card is unavailable")
	}

	if index == 0 {
		p.Hand = p.Hand[1:]
	} else {
		p.Hand = p.Hand[:1]
	}

	return nil
}

func (p *Player) AddHand(card []Card) {
	p.Hand = append(p.Hand, card...)
	// log.Printf("%v hand %v", p.Name, p.Hand)
}
