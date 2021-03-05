package models

// Room :nodoc:
type Room struct {
	RoomID      string   `json:"id_room,omitempty"`
	Capacity    int32    `json:"capacity,omitempty"`
	HostID      string   `json:"id_host,omitempty"`
	IsStarted   bool     `json:"is_started,omitempty"`
	IsClockwise bool     `json:"is_clockwise,omitempty"`
	Players     []Player `json:"players,omitempty"`
	Deck        []Card   `json:"deck,omitempty"`
}
