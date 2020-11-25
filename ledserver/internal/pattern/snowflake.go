package pattern

import (
	"context"
	"math/rand"
)

type Snowflake struct {
	ChanceSpawn float64
}

func (sr *Snowflake) Pattern(ctx context.Context, frameC chan []uint32, numLeds int) {
	colors := make([]uint32, numLeds)
	for {
		// 20% chance every frame to spawn a new snowflake
		if rand.Float64() < sr.ChanceSpawn {
			colors[rand.Intn(numLeds)] = 0xffffff
		}
		for i, color := range colors {
			colors[i] = color
			if colors[i] < 0x080808 {
				colors[i] = 0
			} else {
				colors[i] -= 0x020202
			}
		}
		if WriteOrCancel(ctx, frameC, colors) {
			return
		}
	}
}
