package stealth

import (
	"time"

	"github.com/go-rod/rod"
)

// HoverElement simulates human hover
func HoverElement(el *rod.Element) {
	el.MustHover()
	time.Sleep(800 * time.Millisecond)
}
