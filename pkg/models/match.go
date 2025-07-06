package models

import (
	"time"
)

// Match status constants
const (
	StatusWaiting = "waiting"
	StatusReady   = "ready"
)

// Match represents a game match in the matchmaking system
type Match struct {
	ID        string   `json:"match_id"`
	Players   []string `json:"players"`
	Status    string   `json:"status"`
	UpdatedAt int64    `json:"updated_at"`
	ReadyAt   int64    `json:"ready_at,omitempty"`
}

// NewMatch creates a new match with the given player
func NewMatch(playerID string) *Match {
	now := time.Now().Unix()
	return &Match{
		ID:        generateMatchID(),
		Players:   []string{playerID},
		Status:    StatusWaiting,
		UpdatedAt: now,
	}
}

// AddPlayer adds a player to the match
func (m *Match) AddPlayer(playerID string) {
	m.Players = append(m.Players, playerID)
	m.UpdatedAt = time.Now().Unix()
	
	// If we've reached 3 players, the match is ready
	if len(m.Players) == 3 {
		m.Status = StatusReady
		m.ReadyAt = m.UpdatedAt
	}
}

// IsReady returns true if the match is ready to start
func (m *Match) IsReady() bool {
	return m.Status == StatusReady
}

// generateMatchID creates a unique match ID
func generateMatchID() string {
	return "match_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

// randomString generates a random string of the specified length
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(1 * time.Nanosecond) // Ensure uniqueness
	}
	return string(b)
} 