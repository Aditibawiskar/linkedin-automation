package storage

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type Entry struct {
	ProfileURL string    `json:"profile_url"`
	Action     string    `json:"action"`
	Timestamp  time.Time `json:"timestamp"`
}

var (
	historyFile = "history.json"
	mutex       sync.Mutex
	History     = make(map[string]Entry)
)

func LoadHistory() {
	mutex.Lock()
	defer mutex.Unlock()

	data, err := os.ReadFile(historyFile)
	if err != nil {
		saveToDiskLocked()
		return
	}

	var entries []Entry
	if json.Unmarshal(data, &entries) == nil {
		for _, e := range entries {
			History[e.ProfileURL] = e
		}
	}
}

func IsInvited(url string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	_, exists := History[url]
	return exists
}

func AddInvited(url string) {
	mutex.Lock()
	defer mutex.Unlock()
	entry := Entry{ProfileURL: url, Action: "invited", Timestamp: time.Now()}
	History[url] = entry
	saveToDiskLocked()
}

func saveToDiskLocked() {
	var entries []Entry
	for _, e := range History {
		entries = append(entries, e)
	}
	if entries == nil {
		entries = []Entry{}
	}
	data, _ := json.MarshalIndent(entries, "", "  ")
	os.WriteFile(historyFile, data, 0644)
}
