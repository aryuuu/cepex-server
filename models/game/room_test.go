package game

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
// func ok(tb testing.TB, err error) {
// 	if err != nil {
// 		_, file, line, _ := runtime.Caller(1)
// 		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
// 		tb.FailNow()
// 	}
// }

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\texp: %#v\n\tgot: %#v\033[39m\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestNewRoom(t *testing.T) {
	room := NewRoom("1", "fatt", 2)
	players := []*Player{}

	assert(t, room != nil, "Room should not be nil", room)
	equals(t, "1", room.RoomID)
	equals(t, "fatt", room.HostID)
	equals(t, 2, room.Capacity)
	equals(t, false, room.IsClockwise)
	equals(t, false, room.IsStarted)
	equals(t, 0, room.Count)
	equals(t, players, room.Players)
}

func TestAddPlayer(t *testing.T) {
	player1 := NewPlayer("player1", "")
	player2 := NewPlayer("player2", "")

	room := NewRoom("1", player1.PlayerID, 2)
	room.AddPlayer(player1)
	room.AddPlayer(player2)
	equals(t, 2, len(room.Players))
	equals(t, player1, room.Players[0])
	equals(t, player2, room.Players[1])
}

func TestStartGame(t *testing.T) {
	player1 := NewPlayer("player1", "")
	player2 := NewPlayer("player2", "")
	room := NewRoom("1", "fatt", 2)
	room.AddPlayer(player1)
	room.AddPlayer(player2)
	turnID := room.StartGame()

	equals(t, true, room.IsStarted)
	equals(t, 2, len(player1.Hand))
	equals(t, 2, len(player2.Hand))
	assert(t, turnID == player1.PlayerID || turnID == player2.PlayerID, "Turn ID should be equal to one of players ID")
}

func TestEndGame(t *testing.T) {
	player1 := NewPlayer("player1", "")
	player2 := NewPlayer("player2", "")
	room := NewRoom("1", player1.PlayerID, 2)
	room.AddPlayer(player1)
	room.AddPlayer(player2)
	room.StartGame()

	room.EndGame()

	equals(t, false, room.IsStarted)
	equals(t, false, room.IsClockwise)
	equals(t, 0, room.Count)

	emptyHand := []Card{}
	equals(t, emptyHand, player1.Hand)
	equals(t, emptyHand, player2.Hand)
}
