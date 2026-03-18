package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// UserState stores per-user persistent data: selected role, love score, and
// the last N conversation turns (kept as DeepSeek-compatible Message pairs).
type UserState struct {
	Role    string    `json:"role"`
	Love    int       `json:"love"`
	History []Message `json:"history"`
}

// store is an in-process cache backed by one JSON file per user.
type store struct {
	mu  sync.Mutex
	dir string
}

var defaultStore *store

// InitStore initialises the file-backed store.  dir is the directory used to
// hold per-user JSON files.  It is created if it does not exist.
func InitStore(dir string) error {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("store: mkdir %s: %w", dir, err)
	}
	defaultStore = &store{dir: dir}
	return nil
}

func filePath(openID string) string {
	// openID is already a safe alphanumeric string from WeChat, but sanitise
	// anyway so we never escape the data directory.
	safe := filepath.Base(openID)
	return filepath.Join(defaultStore.dir, safe+".json")
}

// load reads the state for openID without locking (caller must hold the mutex).
func (st *store) load(openID string) (UserState, error) {
	data, err := os.ReadFile(filePath(openID))
	if errors.Is(err, os.ErrNotExist) {
		return UserState{}, nil
	}
	if err != nil {
		return UserState{}, fmt.Errorf("store: read %s: %w", openID, err)
	}
	var s UserState
	if err := json.Unmarshal(data, &s); err != nil {
		return UserState{}, fmt.Errorf("store: unmarshal %s: %w", openID, err)
	}
	return s, nil
}

// save writes the state for openID without locking (caller must hold the mutex).
func (st *store) save(openID string, s UserState) error {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("store: marshal %s: %w", openID, err)
	}
	tmp := filePath(openID) + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("store: write %s: %w", openID, err)
	}
	if err := os.Rename(tmp, filePath(openID)); err != nil {
		return fmt.Errorf("store: rename %s: %w", openID, err)
	}
	return nil
}

// Load reads the state for openID.  Returns a zero-value UserState if the
// file does not exist yet (new user).
func Load(openID string) (UserState, error) {
	defaultStore.mu.Lock()
	defer defaultStore.mu.Unlock()
	return defaultStore.load(openID)
}

// Save writes the state for openID atomically.
func Save(openID string, s UserState) error {
	defaultStore.mu.Lock()
	defer defaultStore.mu.Unlock()
	return defaultStore.save(openID, s)
}

// ClearHistory removes conversation history for openID while preserving the
// role and love score.  The read-modify-write is done under the mutex to
// prevent concurrent modifications.
func ClearHistory(openID string) error {
	defaultStore.mu.Lock()
	defer defaultStore.mu.Unlock()

	s, err := defaultStore.load(openID)
	if err != nil {
		return err
	}
	s.History = nil
	return defaultStore.save(openID, s)
}

// maxHistory caps the stored history to the most recent N pairs (user+assistant).
const maxHistory = 20

// AppendHistory appends a user message and the assistant reply, then trims
// the stored history to maxHistory entries and persists the state.
// The read-modify-write is done under the mutex to prevent concurrent modifications.
func AppendHistory(openID string, userMsg, assistantMsg string) error {
	defaultStore.mu.Lock()
	defer defaultStore.mu.Unlock()

	s, err := defaultStore.load(openID)
	if err != nil {
		return err
	}
	s.History = append(s.History,
		Message{Role: "user", Content: userMsg},
		Message{Role: "assistant", Content: assistantMsg},
	)
	if len(s.History) > maxHistory*2 {
		s.History = s.History[len(s.History)-maxHistory*2:]
	}
	return defaultStore.save(openID, s)
}
