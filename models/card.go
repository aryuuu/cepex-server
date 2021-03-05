package models

// Card :nodoc:
type Card struct {
	Rank    uint32 `json:"rank,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}
