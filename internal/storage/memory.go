package storage

import (
	"encoding/json"
	"os"
)

// Define the filename for our database
const filename = "history.json"

// SaveInvite remembers that we invited this person
func SaveInvite(name string) {
	// 1. Load existing history
	history := LoadHistory()
	
	// 2. Add the new person
	history[name] = true
	
	// 3. Save back to the file
	data, _ := json.MarshalIndent(history, "", "  ")
	os.WriteFile(filename, data, 0644)
}

// LoadHistory reads the file to see who we already invited
func LoadHistory() map[string]bool {
	data, err := os.ReadFile(filename)
	if err != nil {
		// If file doesn't exist, return empty map
		return make(map[string]bool)
	}
	
	var history map[string]bool
	json.Unmarshal(data, &history)
	return history
}

// IsInvited checks if this person is already in our history
func IsInvited(name string) bool {
	history := LoadHistory()
	return history[name]
}