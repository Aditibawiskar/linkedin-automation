package human

import (
	"math/rand"
	"time"
	"github.com/go-rod/rod"
)

// RandomSleep waits for a random time between min and max
func RandomSleep(min, max time.Duration) {
	if min >= max {
		time.Sleep(min)
		return
	}
	delta := max - min
	r := rand.Int63n(int64(delta))
	time.Sleep(min + time.Duration(r))
}

// TypeSlowly types text into an element with random delays
func TypeSlowly(el *rod.Element, text string) {
	for _, char := range text {
		el.MustInput(string(char))
		RandomSleep(50*time.Millisecond, 150*time.Millisecond)
	}
}