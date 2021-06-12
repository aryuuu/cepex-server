package game

import (
	"errors"
	"math/rand"
	"time"

	"github.com/aryuuu/cepex-server/utils/common"
	"github.com/gorilla/websocket"
)

type GameUsecase interface {
	Connect(conn *websocket.Conn, roomID string)
	RunSwitch()
}

// Room :nodoc:
type Room struct {
	RoomID      string             `json:"id_room,omitempty"`
	Capacity    int                `json:"capacity,omitempty"`
	HostID      string             `json:"id_host,omitempty"`
	IsStarted   bool               `json:"is_started,omitempty"`
	IsClockwise bool               `json:"is_clockwise,omitempty"`
	Players     []*Player          `json:"players,omitempty"`
	PlayerMap   map[string]*Player `json:"-"`
	Deck        []Card             `json:"-"`
	TurnID      string             `json:"id_turn"`
	Count       int                `json:"count"`
	VoteBallot  map[string]int     `json:"-"`
}

func NewRoom(id, host string, capacity int) *Room {
	return &Room{
		RoomID:      id,
		Capacity:    capacity,
		HostID:      host,
		IsStarted:   false,
		IsClockwise: false,
		Players:     []*Player{},
		PlayerMap:   make(map[string]*Player),
		Deck:        NewDeck(),
		Count:       0,
		VoteBallot:  make(map[string]int),
	}
}

func (r *Room) StartGame() string {
	r.IsStarted = true

	for _, player := range r.Players {
		player.IsAlive = true
		player.Hand = append(player.Hand, r.PickCard(2)...)
	}

	starterIndex := rand.Intn(len(r.Players))
	r.TurnID = r.Players[starterIndex].PlayerID

	return r.TurnID
}

func (r *Room) EndGame() {
	r.Count = 0
	r.IsStarted = false
	r.IsClockwise = false
	r.Deck = NewDeck()

	for _, p := range r.Players {
		p.IsAlive = false
		p.Hand = []Card{}
	}
}

func (r *Room) PickCard(n int) []Card {
	if len(r.Deck) < n {
		return nil
	}

	result := r.Deck[:n]
	r.Deck = r.Deck[n:]

	return result
}

func (r *Room) PutCard(cards []Card) {
	rand.Seed(time.Now().UnixNano())
	randomNumbers := rand.Perm(len(r.Deck))

	for idx, randomNumber := range randomNumbers[:len(cards)] {
		temp := r.Deck[randomNumber]
		r.Deck[randomNumber] = cards[idx]
		r.Deck = append(r.Deck, temp)
	}
}

func (r *Room) PlayCard(playerID string, handIndex int, isAdd bool, targetID string) error {
	player := r.PlayerMap[playerID]
	card, err := player.PlayHand(handIndex)
	if err != nil {
		return err
	}

	factor := 1
	if !isAdd {
		factor = -1
	}

	if !card.IsSpecial() {
		if r.Count+card.Rank <= 100 {
			r.Count += card.Rank
		} else {
			return errors.New("Card is unplayable")
		}
	} else {
		switch card.Rank {
		case 1:
			if r.Count+(factor*1) > 100 || r.Count+(factor*1) < 0 {
				return errors.New("invalid move")
			}
			r.Count += factor * 1
		case 4:
			r.IsClockwise = !r.IsClockwise
		case 7:
			if target := r.PlayerMap[targetID]; target != nil && !target.IsAlive {
				return errors.New("target is dead")
			}
			r.TurnID = targetID
		case 11:
			if r.Count+(factor*10) > 100 || r.Count+(factor*10) < 0 {
				return errors.New("invalid move")
			}
			r.Count += factor * 10
		case 12:
			if r.Count+(factor*20) > 100 || r.Count+(factor*20) < 0 {
				return errors.New("invalid move")
			}
			r.Count += factor * 20
		case 13:
			r.Count = 100
		default:
			break
		}

	}

	player.AddHand(r.PickCard(1))
	r.PutCard([]Card{card})

	return nil
}

func (r *Room) NextPlayer(playerIndex int) string {
	var nextPlayerIndex int
	if r.IsClockwise {
		for {
			nextPlayerIndex = (playerIndex + 1) % len(r.Players)
			playerIndex++

			if player := r.Players[nextPlayerIndex]; player.IsAlive {
				break
			}
		}

	} else {
		for {
			nextPlayerIndex = common.Mod(playerIndex-1, len(r.Players))
			playerIndex--

			if player := r.Players[nextPlayerIndex]; player.IsAlive {
				break
			}
		}
	}

	r.TurnID = r.Players[nextPlayerIndex].PlayerID

	return r.TurnID
}

func (r *Room) NextHost() (newHostID string) {
	for _, p := range r.Players {
		if p.PlayerID != r.HostID {
			r.HostID = p.PlayerID
			newHostID = p.PlayerID
			return
		}
	}

	return
}

func (r *Room) AddPlayer(player *Player) {
	r.PlayerMap[player.PlayerID] = player
	r.Players = append(r.Players, player)
}

func (r *Room) RemovePlayer(playerIndex int) {
	delete(r.PlayerMap, r.Players[playerIndex].PlayerID)
	r.Players = append(r.Players[:playerIndex], r.Players[playerIndex+1:]...)
}

func (r *Room) GetPlayerIndex(playerID string) int {
	playerIndex := -1
	for i, p := range r.Players {
		if p.PlayerID == playerID {
			playerIndex = i
			break
		}
	}

	return playerIndex
}

func (r *Room) GetWinner() (winner string) {
	survivor := []string{}

	for _, p := range r.Players {
		if p.IsAlive {
			survivor = append(survivor, p.PlayerID)
		}
	}

	if len(survivor) == 1 {
		return survivor[0]
	}

	return
}
