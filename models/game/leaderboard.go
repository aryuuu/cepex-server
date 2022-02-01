package game

type Leaderboard struct {
	Items []LeaderboardItem `json:"items"`
}

type LeaderboardItem struct {
	PlayerID  string `json:"id_player,omitempty"`
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Score     int    `json:"score,omitempty"`
}
