package pattern

import (
	"context"
	"math"
	"pnpleds/ledserver/internal"
)

type Breathe struct {
	Colors []uint32
}

// Smooths the transition between colors by "accelerating"
// and "decelerating" near the endpoints of the interpolation.
func transform(t float64) float64 {
	return -math.Cos(t*math.Pi)/2 + 0.5
}

func (b *Breathe) Pattern(ctx context.Context, frameC chan []uint32, numLeds int) {
	leds := make([]uint32, numLeds)
	numColors := len(b.Colors)
	nxt := func(curr int) int {
		curr++
		if curr == numColors {
			curr = 0
		}
		return curr
	}
	currColor := 0
	nextColor := nxt(currColor)
	t := 0.
	for {
		for i := range leds {
			leds[i] = internal.Lerp(
				b.Colors[currColor],
				b.Colors[nextColor],
				transform(t),
			)
		}
		t += 0.01
		if t > 1.0 {
			t = 0.
			currColor = nextColor
			nextColor = nxt(nextColor)
		}
		if WriteOrCancel(ctx, frameC, leds) {
			return
		}
	}
}
