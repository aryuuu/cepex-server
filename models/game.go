package models

import "github.com/gorilla/websocket"

// Card :nodoc:
type Card struct {
	Rank    uint32 `json:"rank,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

// Player :nodoc:
type Player struct {
	PlayerID string `json:"id_player,omitplayer"`
	Name     string `json:"name,omitempty"`
	IsAlive  bool   `json:"is_alive,omitempty"`
	Hand     []Card `json:"hand,omitempty"`
}

// Room :nodoc:
type Room struct {
	RoomID      string            `json:"id_room,omitempty"`
	Capacity    int32             `json:"capacity,omitempty"`
	HostID      string            `json:"id_host,omitempty"`
	IsStarted   bool              `json:"is_started,omitempty"`
	IsClockwise bool              `json:"is_clockwise,omitempty"`
	Players     map[string]Player `json:"players,omitempty"`
	Deck        []Card            `json:"-"`
	Count       int32             `json:"count,omitempty"`
}

type SocketServer struct {
	clients map[uint32]*SocketClient
}

type SocketClient struct {
	ID   uint32
	conn *websocket.Conn
}

type GameUsecase interface {
	SendMessage(connID string, message interface{})
	BroadcastMessage(roomID string, message interface{})
}
