package game

import (
	"math/rand"
	"time"
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

func (c Card) IsSpecial() bool {
	return c.Rank == 1 || c.Rank == 4 || c.Rank == 7 || c.Rank == 11 || c.Rank == 12 || c.Rank == 13
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
