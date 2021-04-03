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
	EventType string        `json:"event_type,omitempty"`
	Success   bool          `json:"success,omitempty"`
	NewRoom   models.Room   `json:"room,omitempty"`
	Hand      []models.Card `json:"hand,omitempty"`
}

type JoinRoomResponse struct {
	EventType string        `json:"event_type,omitempty"`
	Success   bool          `json:"success,omitempty"`
	NewRoom   models.Room   `json:"new_room,omitempty"`
	Hand      []models.Card `json:"hand,omitempty"`
}

type JoinRoomBroadcast struct {
	EventType string        `json:"event_type,omitempty"`
	NewPlayer models.Player `json:"new_player,omitempty"`
}

type LeaveRoomResponse struct {
	EventType string `json:"event_type,omitempty"`
	Success   bool   `json:"success,omitempty"`
}

type LeaveRoomBroadcast struct {
	EventType       string `json:"event_type,omitempty"`
	LeavingPlayerID string `json:"id_leaving_player,omitempty"`
}

type StartGameResponse struct {
	EventType string `json:"event_type,omitempty"`
	Success   bool   `json:"success,omitempty"`
	StarterID string `json:"id_starter,omitempty"`
}

type StartGameBroadcast struct {
	EventType string `json:"event_type,omitempty"`
	StarterID string `json:",omitempty"`
}

type PlayCardResponse struct {
	EventType string `json:"event_type,omitempty"`
	Success   bool   `json:"success,omitempty"`
}

type PlayCardBroadcast struct {
	EventType string `json:"event_type,omitempty"`
	Count     int32  `json:"count,omitempty"`
}

type TurnBroadcast struct {
	EventType string `json:"event_type,omitempty"`
	PlayerID  string `json:"id_player,omitempty"`
}

type MessageBroadcast struct {
	EventType string `json:"event_type,omitempty"`
	Sender    string `json:"sender,emitempty"`
	Message   string `json:"message,omitempty"`
}

func NewCreateRoomResponse(success bool, roomID string, host models.Player, hand []models.Card) CreateRoomResponse {
	players := make(map[string]models.Player)
	players[host.Name] = host

	result := CreateRoomResponse{
		EventType: "create-room",
		Success:   success,
		NewRoom: models.Room{
			RoomID:      roomID,
			Capacity:    4,
			HostID:      host.PlayerID,
			IsStarted:   false,
			IsClockwise: false,
			Players:     players,
			Count:       0,
		},
		Hand: hand,
	}

	return result
}

func NewJoinRoomResponse(success bool, room models.Room) JoinRoomResponse {
	result := JoinRoomResponse{
		EventType: "join-room",
		Success:   success,
		NewRoom:   room,
	}

	return result
}

func NewJoinRoomBroadcast(player models.Player) JoinRoomBroadcast {
	result := JoinRoomBroadcast{
		EventType: "join-room-broadcast",
		NewPlayer: models.Player{},
	}

	return result
}

func NewLeaveRoomResponse(success bool) LeaveRoomResponse {
	result := LeaveRoomResponse{
		EventType: "leave-room",
		Success:   success,
	}

	return result
}

func NewLeaveRoomBroadcast(playerID string) LeaveRoomBroadcast {
	result := LeaveRoomBroadcast{
		EventType:       "leave-room-broadcast",
		LeavingPlayerID: playerID,
	}

	return result
}

func NewMessageBroadcast(message, sender string) MessageBroadcast {
	result := MessageBroadcast{
		EventType: "message-broadcast",
		Message:   message,
		Sender:    sender,
	}

	return result
}
