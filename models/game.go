package models

import (
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

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
	RoomID      string   `json:"id_room,omitempty"`
	Capacity    int32    `json:"capacity,omitempty"`
	HostID      string   `json:"id_host,omitempty"`
	IsStarted   bool     `json:"is_started,omitempty"`
	IsClockwise bool     `json:"is_clockwise,omitempty"`
	Players     []Player `json:"players,omitempty"`
	Deck        []Card   `json:"-"`
	Count       int32    `json:"count"`
}

type SocketServer struct {
	clients map[uint32]*SocketClient
}

type SocketClient struct {
	ID   uint32
	conn *websocket.Conn
}

type GameUsecase interface {
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

func (r *Room) AddPlayer(player Player) {
	r.Players = append(r.Players, player)
}

func (r *Room) RemovePlayer(playerIndex int) {
	r.Players = append(r.Players[:playerIndex], r.Players[playerIndex+1:]...)
}
