package human

import (
	"github.com/go-rod/rod"
	"math/rand"
	"time"
)

func RandomSleep(min, max time.Duration) {
	if min >= max {
		time.Sleep(min)
		return
	}
	delta := max - min
	r := rand.Int63n(int64(delta))
	time.Sleep(min + time.Duration(r))
}

func TypeSlowly(el *rod.Element, text string) {
	for _, char := range text {
		el.MustInput(string(char))
		RandomSleep(50*time.Millisecond, 150*time.Millisecond)
	}
}
