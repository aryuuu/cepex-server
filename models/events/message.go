package events

import "github.com/aryuuu/cepex-server/models"

type GameRequest struct {
	EventType  string `json:"event_type,omitempty"`
	ClientName string `json:"client_name"`
	AvatarURL  string `json:"avatar_url"`
	Message    string `json:"message,omitempty"`
	HandIndex  int    `json:"hand_index,omitempty"`
	IsAdd      bool   `json:"is_add,omitempty"`
	PlayerID   string `json:"id_player,omitempty"`
}

type GameResponse struct {
	EventType string          `json:"event_type,omitempty"`
	Players   []models.Player `json:"players,omitempty"`
}

type CreateRoomResponse struct {
	EventType string      `json:"event_type,omitempty"`
	Success   bool        `json:"success,omitempty"`
	NewRoom   models.Room `json:"room,omitempty"`
	// Hand      []models.Card `json:"hand"`
}

type JoinRoomResponse struct {
	EventType string      `json:"event_type,omitempty"`
	Success   bool        `json:"success"`
	NewRoom   models.Room `json:"new_room,omitempty"`
	// Hand      []models.Card `json:"hand"`
}

type JoinRoomBroadcast struct {
	EventType string         `json:"event_type,omitempty"`
	NewPlayer *models.Player `json:"new_player,omitempty"`
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
	EventType    string `json:"event_type,omitempty"`
	StarterIndex int    `json:"starter_idx"`
}

type EndGameBroadcast struct {
	EventType string `json:"event_type"`
	WinnerID  string `json:"id_winner,omitempty"`
}

type InitialHandResponse struct {
	EventType string        `json:"event_type"`
	NewHand   []models.Card `json:"new_hand"`
}

type PlayCardResponse struct {
	EventType string        `json:"event_type,omitempty"`
	Success   bool          `json:"success"`
	IsUpdate  bool          `json:"is_update"`
	NewHand   []models.Card `json:"new_hand"`
	// HandIndex int         `json:"hand_index,omitempty"`
	// DrawnCard models.Card `json:"drawn_card,omitempty"`
}

type PlayCardBroadcast struct {
	EventType       string      `json:"event_type,omitempty"`
	Card            models.Card `json:"card"`
	Count           int         `json:"count,omitempty"`
	IsClockwise     bool        `json:"is_clockwise"`
	NextPlayerIndex int         `json:"next_player_idx"`
}

type TurnBroadcast struct {
	EventType string `json:"event_type,omitempty"`
	PlayerID  string `json:"id_player,omitempty"`
}

type MessageBroadcast struct {
	EventType string `json:"event_type,omitempty"`
	Sender    string `json:"sender,omitempty"`
	Message   string `json:"message,omitempty"`
}

type NotificationBroadcast struct {
	EventType string `json:"event_type,omitempty"`
	Message   string `json:"message,omitempty"`
}

type DeadPlayerBroadcast struct {
	EventType    string `json:"event_type"`
	DeadPlayerID string `json:"id_dead_player"`
}

func NewCreateRoomResponse(success bool, roomID string, host *models.Player) CreateRoomResponse {
	players := []*models.Player{host}

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
	}

	return result
}

func NewJoinRoomResponse(success bool, room *models.Room) JoinRoomResponse {
	result := JoinRoomResponse{
		EventType: "join-room",
		Success:   success,
		NewRoom:   *room,
	}

	return result
}

func NewJoinRoomBroadcast(player *models.Player) JoinRoomBroadcast {
	result := JoinRoomBroadcast{
		EventType: "join-room-broadcast",
		NewPlayer: player,
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

func NewNotificationBroadcast(message string) NotificationBroadcast {
	result := NotificationBroadcast{
		EventType: "notification-broadcast",
		Message:   message,
	}

	return result
}

func NewStartGameResponse(success bool) StartGameResponse {
	result := StartGameResponse{
		EventType: "start-game",
		Success:   success,
	}

	return result
}

func NewStartGameBroadcast(starterIndex int) StartGameBroadcast {
	result := StartGameBroadcast{
		EventType:    "start-game-broadcast",
		StarterIndex: starterIndex,
	}

	return result
}

func NewEndGameBroadcast(winnerID string) EndGameBroadcast {
	result := EndGameBroadcast{
		EventType: "end-game-broadcast",
		WinnerID:  winnerID,
	}

	return result
}

func NewInitialHandResponse(hand []models.Card) InitialHandResponse {
	result := InitialHandResponse{
		EventType: "initial-hand",
		NewHand:   hand,
	}

	return result
}

func NewPlayCardResponse(success bool, newHand []models.Card) PlayCardResponse {
	result := PlayCardResponse{
		EventType: "play-card",
		Success:   success,
		NewHand:   newHand,
		IsUpdate:  newHand != nil,
	}

	return result
}

func NewPlayCardBroadcast(card models.Card, count int, isClockwise bool, nextPlayerIdx int) PlayCardBroadcast {
	result := PlayCardBroadcast{
		EventType:       "play-card-broadcast",
		Card:            card,
		Count:           count,
		IsClockwise:     isClockwise,
		NextPlayerIndex: nextPlayerIdx,
	}

	return result
}

func NewTurnBroadcast(playerID string) TurnBroadcast {
	result := TurnBroadcast{
		EventType: "turn-broadcast",
		PlayerID:  playerID,
	}

	return result
}

func NewDeadPlayerBroadcast(playerID string) DeadPlayerBroadcast {
	result := DeadPlayerBroadcast{
		EventType:    "dead-player",
		DeadPlayerID: playerID,
	}

	return result
}
