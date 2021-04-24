package models

import (
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GameUsecase interface {
	Connect(conn *websocket.Conn, roomID string)
}

// Card :nodoc:
type Card struct {
	// 0 = diamond
	// 1 = club
	// 2 = heart
	// 3 = spade
	Pattern int `json:"pattern"`
	Rank    int `json:"rank"`
}

// Player :nodoc:
type Player struct {
	PlayerID  string `json:"id_player,omitplayer"`
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatar_url"`
	IsAlive   bool   `json:"is_alive"`
	Hand      []Card `json:"-"`
}

// Room :nodoc:
type Room struct {
	RoomID      string    `json:"id_room,omitempty"`
	Capacity    int       `json:"capacity,omitempty"`
	HostID      string    `json:"id_host,omitempty"`
	IsStarted   bool      `json:"is_started,omitempty"`
	IsClockwise bool      `json:"is_clockwise,omitempty"`
	Players     []*Player `json:"players,omitempty"`
	Deck        []Card    `json:"-"`
	TurnID      string    `json:"id_turn"`
	Count       int       `json:"count"`
}

type SocketServer struct {
	clients map[uint32]*SocketClient
}

type SocketClient struct {
	ID   uint32
	conn *websocket.Conn
}

func NewPlayer(name, avatarUrl string) *Player {
	return &Player{
		Name:      name,
		AvatarURL: avatarUrl,
		PlayerID:  uuid.NewString(),
		IsAlive:   true,
		Hand:      []Card{},
	}
}

func NewRoom(id, host string, capacity int) *Room {
	return &Room{
		RoomID:      id,
		Capacity:    capacity,
		HostID:      host,
		IsStarted:   false,
		IsClockwise: false,
		Players:     []*Player{},
		Deck:        NewDeck(),
		Count:       0,
	}
}

func NewDeck() []Card {
	totalCard := 52
	result := make([]Card, totalCard)

	for pattern := 0; pattern < 4; pattern++ {
		for rank := 1; rank < 14; rank++ {
			result[(13*pattern)+rank-1] = Card{
				Rank:    rank,
				Pattern: pattern,
			}
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })

	return result
}

func (c Card) IsSpecial() bool {
	return c.Rank == 1 || c.Rank == 4 || c.Rank == 7 || c.Rank == 11 || c.Rank == 12 || c.Rank == 13
}

func (r *Room) StartGame() bool {
	r.IsStarted = true

	for _, player := range r.Players {
		player.IsAlive = true
	}

	return true
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

func (r *Room) PlayCard(card Card, isAdd bool) bool {
	factor := 1
	if !isAdd {
		factor = -1
	}

	if !card.IsSpecial() {
		if r.Count+card.Rank <= 100 {
			r.Count += card.Rank
		} else {
			return false
		}
	} else {
		switch card.Rank {
		case 1:
			r.Count += factor * 1
			break
		case 4:
			r.IsClockwise = !r.IsClockwise
			break
		case 7:
			break
		case 11:
			r.Count += factor * 10
		case 12:
			r.Count += factor * 20
		case 13:
			r.Count = 100
		default:
			break
		}

	}
	return true
}

func (r *Room) AddPlayer(player *Player) {
	r.Players = append(r.Players, player)
}

func (r *Room) RemovePlayer(playerIndex int) {
	r.Players = append(r.Players[:playerIndex], r.Players[playerIndex+1:]...)
}

func (p *Player) PlayHand(index int) bool {
	if index >= len(p.Hand) {
		return false
	}

	if index == 0 {
		p.Hand = p.Hand[1:]
	} else {
		p.Hand = p.Hand[:1]
	}

	return true
}

func (p *Player) AddHand(card []Card) {
	log.Printf("player hand")
	p.Hand = append(p.Hand, card...)
}
