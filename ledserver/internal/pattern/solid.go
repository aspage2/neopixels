package pattern

import "context"

type Solid struct {
	Color uint32
}

func (s *Solid) Pattern(ctx context.Context, frameC chan []uint32, numLeds int) {
	leds := make([]uint32, numLeds)
	for i := range leds {
		leds[i] = s.Color
	}

	for {
		if WriteOrCancel(ctx, frameC, leds) {
			return
		}
	}

}
