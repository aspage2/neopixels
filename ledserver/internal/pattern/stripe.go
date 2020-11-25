package pattern

import "context"

// Stripe collects config related to a "color stripe" effect"
type Stripe struct {
	// StripeSize is the number of LEDs per "color stripe"
	StripeSize int

	// The colors to cycle through for each stripe
	Colors []uint32
}

func (sr *Stripe) Pattern(ctx context.Context, frameC chan []uint32, numLeds int) {
	N := len(sr.Colors)
	cycleLen := N * sr.StripeSize
	leds := make([]uint32, numLeds)
	clk := 0
	for {
		for i := range leds {
			ind := ((i + clk) / sr.StripeSize) % N
			leds[i] = sr.Colors[ind]
		}
		clk = (clk + 1) % cycleLen

		if WriteOrCancel(ctx, frameC, leds) {
			return
		}
	}
}
