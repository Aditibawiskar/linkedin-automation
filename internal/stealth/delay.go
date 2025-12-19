package stealth

import (
	"math/rand"
	"time"
)

// SleepRandom pauses execution for a human-like random duration
func SleepRandom(minMs, maxMs int) {
	rand.Seed(time.Now().UnixNano())
	delay := rand.Intn(maxMs-minMs) + minMs
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
