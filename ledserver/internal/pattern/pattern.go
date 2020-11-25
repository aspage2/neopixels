package pattern

import "context"

// A Pattern generates a sequence of "frames" to display on the LED strip
type Pattern interface {
	Pattern(context.Context, chan []uint32, int)
}

func WriteOrCancel(ctx context.Context, frameC chan []uint32, frame []uint32) (shouldExit bool) {
	select {
	case <-ctx.Done():
		shouldExit = true
	case frameC <- frame:
		shouldExit = false
	}
	return
}
