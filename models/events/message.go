package events

import (
	"github.com/aryuuu/cepex-server/models/game"
	"github.com/gorilla/websocket"
)

type SocketEvent struct {
	EventType string          `json:"event_type"`
	RoomID    string          `json:"id_room"`
	Conn      *websocket.Conn `json:"conn"`
	Message   interface{}     `json:"message"`
}

type GameRequest struct {
	EventType  string `json:"event_type,omitempty"`
	ClientName string `json:"client_name"`
	AvatarURL  string `json:"avatar_url"`
	Message    string `json:"message,omitempty"`
	HandIndex  int    `json:"hand_index,omitempty"`
	IsAdd      bool   `json:"is_add,omitempty"`
	PlayerID   string `json:"id_player,omitempty"`
	IsDiscard  bool   `json:"is_discard"`
}

type GameResponse struct {
	EventType string        `json:"event_type,omitempty"`
	Players   []game.Player `json:"players,omitempty"`
}

type CreateRoomResponse struct {
	EventType string    `json:"event_type,omitempty"`
	Success   bool      `json:"success,omitempty"`
	NewRoom   game.Room `json:"room,omitempty"`
	Detail    string    `json:"detail,omitempty"`
	// Hand      []game.Card `json:"hand"`
}

type JoinRoomResponse struct {
	EventType string    `json:"event_type,omitempty"`
	Success   bool      `json:"success"`
	NewRoom   game.Room `json:"new_room,omitempty"`
	Detail    string    `json:"detail,omitempty"`
	// Hand      []game.Card `json:"hand"`
}

type JoinRoomBroadcast struct {
	EventType string       `json:"event_type,omitempty"`
	NewPlayer *game.Player `json:"new_player,omitempty"`
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

type LeaderboardResponse struct {
	EventType   string           `json:"event_type,omitempty"`
	Leaderboard game.Leaderboard `json:"leaderboard,omitempty"`
}

type StartGameBroadcast struct {
	EventType string `json:"event_type,omitempty"`
	StarterID string `json:"id_starter"`
}

type EndGameBroadcast struct {
	EventType   string `json:"event_type"`
	WinnerID    string `json:"id_winner,omitempty"`
	WinnerScore int    `json:"winner_score"`
}

type InitialHandResponse struct {
	EventType string      `json:"event_type"`
	NewHand   []game.Card `json:"new_hand"`
}

type PlayCardResponse struct {
	EventType string      `json:"event_type,omitempty"`
	Success   bool        `json:"success"`
	IsUpdate  bool        `json:"is_update"`
	NewHand   []game.Card `json:"new_hand"`
	Message   string      `json:"message,omitempty"`
	Status    int         `json:"status"`
	HandIndex int         `json:"hand_index"`
	// status code list
	// 0: success no prob
	// 1: unplayable card
	// 2: discard success
	// 3: other error
}

type PlayCardBroadcast struct {
	EventType    string    `json:"event_type,omitempty"`
	Card         game.Card `json:"card"`
	Count        int       `json:"count,omitempty"`
	IsClockwise  bool      `json:"is_clockwise"`
	NextPlayerID string    `json:"id_next_player"`
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

type ChangeHostBroadcast struct {
	EventType string `json:"event_type"`
	NewHostID string `json:"id_new_host"`
}

type VoteKickPlayerResponse struct {
	EventType string `json:"event_type"`
	Success   bool   `json:"success"`
}

type VoteKickPlayerBroadcast struct {
	EventType  string `json:"event_type"`
	TargetID   string `json:"id_target"`
	IssuerName string `json:"issuer_name"`
}

func NewUnicastEvent(roomID string, conn *websocket.Conn, message interface{}) SocketEvent {
	return SocketEvent{
		EventType: "unicast",
		RoomID:    roomID,
		Conn:      conn,
		Message:   message,
	}
}

func NewBroadcastEvent(roomID string, message interface{}) SocketEvent {
	return SocketEvent{
		EventType: "broadcast",
		RoomID:    roomID,
		Message:   message,
	}
}

func NewSocketEvent(eventType, roomID string, conn *websocket.Conn, message interface{}) SocketEvent {
	return SocketEvent{
		EventType: eventType,
		RoomID:    roomID,
		Message:   message,
	}
}

func NewCreateRoomResponse(success bool, roomID string, host *game.Player, detail string) CreateRoomResponse {
	players := []*game.Player{host}

	result := CreateRoomResponse{
		EventType: "create-room",
		Success:   success,
		NewRoom: game.Room{
			RoomID:      roomID,
			Capacity:    4,
			HostID:      host.PlayerID,
			IsStarted:   false,
			IsClockwise: false,
			Players:     players,
			Count:       0,
		},
		Detail: detail,
	}

	return result
}

func NewJoinRoomResponse(success bool, room *game.Room, detail string) JoinRoomResponse {
	result := JoinRoomResponse{
		EventType: "join-room",
		Success:   success,
		NewRoom:   *room,
		Detail:    detail,
	}

	return result
}

func NewJoinRoomBroadcast(player *game.Player) JoinRoomBroadcast {
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

func NewStartGameBroadcast(starterID string) StartGameBroadcast {
	result := StartGameBroadcast{
		EventType: "start-game-broadcast",
		StarterID: starterID,
	}

	return result
}

func NewEndGameBroadcast(winner *game.Player) EndGameBroadcast {
	result := EndGameBroadcast{
		EventType:   "end-game-broadcast",
		WinnerID:    winner.PlayerID,
		WinnerScore: winner.Score,
	}

	return result
}

func NewInitialHandResponse(hand []game.Card) InitialHandResponse {
	result := InitialHandResponse{
		EventType: "initial-hand",
		NewHand:   hand,
	}

	return result
}

func NewPlayCardResponse(success bool, newHand []game.Card, status int, message string) PlayCardResponse {
	result := PlayCardResponse{
		EventType: "play-card",
		Success:   success,
		NewHand:   newHand,
		IsUpdate:  newHand != nil,
		Status:    status,
		Message:   message,
	}

	return result
}

func NewPlayCardBroadcast(card game.Card, count int, isClockwise bool, nextPlayerID string) PlayCardBroadcast {
	result := PlayCardBroadcast{
		EventType:    "play-card-broadcast",
		Card:         card,
		Count:        count,
		IsClockwise:  isClockwise,
		NextPlayerID: nextPlayerID,
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

func NewChangeHostBroadcast(hostID string) ChangeHostBroadcast {
	return ChangeHostBroadcast{
		EventType: "change-host",
		NewHostID: hostID,
	}
}

func NewVoteKickPlayerResponse(success bool) VoteKickPlayerResponse {
	return VoteKickPlayerResponse{
		EventType: "vote-kick",
		Success:   success,
	}
}

func NewVoteKickPlayerBroadcast(targetID, issuerName string) VoteKickPlayerBroadcast {
	return VoteKickPlayerBroadcast{
		EventType:  "vote-kick-broadcast",
		TargetID:   targetID,
		IssuerName: issuerName,
	}
}
