package models

// Player :nodoc:
type Player struct {
	PlayerID string `json:"id_player,omitplayer"`
	IsAlive  bool   `json:"is_alive,omitempty"`
	Hand     []Card `json:"hand,omitempty"`
}
