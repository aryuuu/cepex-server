package events

import "github.com/aryuuu/cepex-server/models"

type GameRequest struct {
	EventType string `json:"event_type,omitempty"`
	Message   string `json:"message,omitempty"`
}

type GameResponse struct {
	EventType string          `json:"event_type,omitempty"`
	Players   []models.Player `json:"players,omitempty"`
}

type CreateRoomResponse struct {
	EventType string      `json:"event_type,omitempty"`
	Success   bool        `json:"success,omitempty"`
	NewRoom   models.Room `json:"room,omitempty"`
}

type JoinRoomResponse struct {
	EventType string      `json:",omitempty"`
	Success   bool        `json:",omitempty"`
	NewRoom   models.Room `json:",omitempty"`
}

type JoinRoomBroadcast struct {
	EventType string        `json:",omitempty"`
	NewPlayer models.Player `json:",omitempty"`
}

type LeaveRoomResponse struct {
	EventType string `json:",omitempty"`
	Success   bool   `json:",omitempty"`
}

type LeaveRoomBroadcast struct {
	EventType       string `json:",omitempty"`
	LeavingPlayerID string `json:",omitempty"`
}

type StartGameResponse struct {
	EventType string `json:",omitempty"`
	Success   bool   `json:",omitempty"`
	StarterID string `json:",omitempty"`
}

type StartGameBroadcast struct {
	EventType string `json:",omitempty"`
	StarterID string `json:",omitempty"`
}

type PlayCardResponse struct {
	EventType string `json:",omitempty"`
	Success   bool   `json:",omitempty"`
}

type PlayCardBroadcast struct {
	EventType string `json:",omitempty"`
	Count     int32  `json:",omitempty"`
}

type TurnBroadcast struct {
	EventType string `json:",omitempty"`
	PlayerID  string `json:",omitempty"`
}

/* events
create room
join room
leave room
start game
play card
pass

*/
