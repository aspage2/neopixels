package internal

import (
	"github.com/Jon-Bright/ledctl/pixarray"
)

type StripOptions struct {
	NumPixels    int
	Order        int
	OscFrequency uint
	DMAChannel   int
	PWMPins      []int
}

type option func(*StripOptions)

func WithPixelOrder(order int) option {
	return func(opts *StripOptions) {
		opts.Order = order
	}
}

func WithOscFreq(freq uint) option {
	return func(opts *StripOptions) {
		opts.OscFrequency = freq
	}
}

func WithDMAChannel(channel int) option {
	return func(opts *StripOptions) {
		opts.DMAChannel = channel
	}
}

func NewStrip(numPixels int, options ...option) (pixarray.LEDStrip, error) {
	opts := StripOptions{
		NumPixels:    numPixels,
		Order:        pixarray.RGB,
		OscFrequency: 800000,
		DMAChannel:   10,
		PWMPins:      []int{18},
	}

	for _, o := range options {
		o(&opts)
	}

	return pixarray.NewWS281x(
		opts.NumPixels,
		3,
		opts.Order,
		uint(opts.OscFrequency),
		opts.DMAChannel,
		opts.PWMPins,
	)
}
