package matchmaker

import (
	"errors"
	"sync"

	"matchmaking-httpapi/pkg/models"
)

var (
	// ErrMatchNotFound is returned when a match with the specified ID is not found
	ErrMatchNotFound = errors.New("match not found")
)

// Matchmaker handles the matchmaking process for the game
type Matchmaker struct {
	matches       map[string]*models.Match // Map of match ID to match
	waitingPlayers []string              // Queue of players waiting to be matched
	pendingMatches map[string]*models.Match // Map of incomplete matches
	mu            sync.RWMutex            // Mutex for concurrent access
}

// NewMatchmaker creates a new matchmaker instance
func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		matches:        make(map[string]*models.Match),
		waitingPlayers: make([]string, 0),
		pendingMatches: make(map[string]*models.Match),
		mu:             sync.RWMutex{},
	}
}

// AddPlayer adds a player to the matchmaking queue and tries to form a match
func (m *Matchmaker) AddPlayer(playerID string) (*models.Match, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Add player to the waiting queue
	m.waitingPlayers = append(m.waitingPlayers, playerID)

	// Try to form matches
	match := m.formMatches(playerID)
	
	return match, nil
}

// formMatches attempts to form matches with waiting players
// It returns the match that the provided playerID was added to
func (m *Matchmaker) formMatches(playerID string) *models.Match {
	var playerMatch *models.Match

	// If we have at least 3 players waiting, form a match
	for len(m.waitingPlayers) >= 3 {
		// Take the first 3 players from the queue
		players := m.waitingPlayers[:3]
		m.waitingPlayers = m.waitingPlayers[3:]

		// Create a new match with these players
		match := models.NewMatch(players[0])
		match.AddPlayer(players[1])
		match.AddPlayer(players[2])
		
		// Store the match
		m.matches[match.ID] = match
		
		// If this match contains our player, remember it
		for _, p := range players {
			if p == playerID {
				playerMatch = match
				break
			}
		}
	}

	// If our player wasn't matched yet, create a pending match
	if playerMatch == nil {
		// Create a new pending match
		match := models.NewMatch(playerID)
		m.pendingMatches[match.ID] = match
		playerMatch = match
	}

	return playerMatch
}

// GetMatch returns the match with the specified ID
func (m *Matchmaker) GetMatch(matchID string) (*models.Match, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Check complete matches
	if match, ok := m.matches[matchID]; ok {
		return match, nil
	}

	// Check pending matches
	if match, ok := m.pendingMatches[matchID]; ok {
		return match, nil
	}

	return nil, ErrMatchNotFound
} 