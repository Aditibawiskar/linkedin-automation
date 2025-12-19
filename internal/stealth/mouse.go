package stealth

import (
	"math"
	"math/rand"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

func MoveMouseHuman(page *rod.Page, startX, startY, endX, endY float64) {
	steps := rand.Intn(20) + 30

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)

		curve := math.Sin(t * math.Pi)
		x := startX + (endX-startX)*t
		y := startY + (endY-startY)*t + curve*rand.Float64()*8

		page.Mouse.MoveTo(proto.Point{
			X: x,
			Y: y,
		})

		time.Sleep(time.Duration(rand.Intn(20)+10) * time.Millisecond)
	}
}
