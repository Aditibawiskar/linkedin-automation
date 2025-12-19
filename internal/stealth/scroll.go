package stealth

import (
	"math/rand"
	"time"

	"github.com/go-rod/rod"
)

// RandomScroll simulates human-like scrolling using mouse wheel
func RandomScroll(page *rod.Page) {
	scrollTimes := rand.Intn(4) + 3

	for i := 0; i < scrollTimes; i++ {
		x := rand.Float64()*400 + 200
		y := rand.Float64()*400 + 200
		delta := rand.Intn(400) + 200

		page.Mouse.Scroll(x, y, delta)

		time.Sleep(time.Duration(rand.Intn(1200)+600) * time.Millisecond)
	}
}
