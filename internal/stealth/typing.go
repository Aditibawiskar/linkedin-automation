package stealth

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
)

// TypeLikeHuman types text with variable delays
func TypeLikeHuman(el *rod.Element, text string) {
	for _, ch := range text {
		el.MustInput(string(ch))
		time.Sleep(time.Duration(rand.Intn(120)+80) * time.Millisecond)
	}
}
