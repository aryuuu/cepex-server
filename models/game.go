package models

import (
	"errors"
	"math/rand"
	"time"

	"github.com/aryuuu/cepex-server/utils/common"
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
	PlayerID  string `json:"id_player,omitempty"`
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatar_url"`
	IsAlive   bool   `json:"is_alive"`
	Hand      []Card `json:"-"`
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
}

// type SocketServer struct {
// 	clients map[uint32]*SocketClient
// }

// type SocketClient struct {
// 	ID   uint32
// 	conn *websocket.Conn
// }

func NewPlayer(name, avatarUrl string) *Player {
	return &Player{
		Name:      name,
		AvatarURL: avatarUrl,
		PlayerID:  uuid.NewString(),
		IsAlive:   false,
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
		PlayerMap:   make(map[string]*Player),
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

func (r *Room) StartGame() int {
	r.IsStarted = true

	for _, player := range r.Players {
		player.IsAlive = true
		player.Hand = append(player.Hand, r.PickCard(2)...)
	}

	starterIndex := rand.Intn(len(r.Players))
	r.TurnID = r.Players[starterIndex].PlayerID

	return starterIndex
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

func (r *Room) PlayCard(playerID string, handIndex int, isAdd bool) error {
	player := r.PlayerMap[playerID]
	card := player.Hand[handIndex]

	if err := player.PlayHand(handIndex); err != nil {
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
			r.Count += factor * 1
		case 4:
			r.IsClockwise = !r.IsClockwise
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

	player.AddHand(r.PickCard(1))
	r.PutCard([]Card{card})

	return nil
}

func (r *Room) NextPlayer(playerIndex int) int {
	var nextPlayerIndex int
	if r.IsClockwise {
		nextPlayerIndex = (playerIndex + 1) % len(r.Players)

	} else {
		nextPlayerIndex = common.Mod(playerIndex-1, len(r.Players))
	}

	r.TurnID = r.Players[nextPlayerIndex].PlayerID

	return nextPlayerIndex
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
