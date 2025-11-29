//go:build arm
package main

import (
	"strings"

	"github.com/Jon-Bright/ledctl/pixarray"
)

type StripOptions struct {
	NumPixels    int
	Order        int
	OscFrequency uint
	DMAChannel   int
	PWMPins      []int
	Brightness   float32
}

type option func(*StripOptions)

func WithPixelOrder(order string) option {
	return func(opts *StripOptions) {
		opts.Order = pixarray.StringOrders[strings.ToUpper(order)]
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

func WithBrightness(b float32) option {
	return func(opts *StripOptions) {
		opts.Brightness = b
	}
}

type LEDStripWithBrightness struct {
	pixarray.LEDStrip
	Brightness float32
}

func (s *LEDStripWithBrightness) SetPixel(i int, p pixarray.Pixel) {
	s.LEDStrip.SetPixel(i, GammaCorrect(Scale(p, s.Brightness)))
}

func NewStrip(numPixels int, options ...option) (pixarray.LEDStrip, error) {
	opts := StripOptions{
		NumPixels:    numPixels,
		Order:        pixarray.RGB,
		OscFrequency: 800000,
		DMAChannel:   10,
		PWMPins:      []int{18},
		Brightness:   0,
	}

	for _, o := range options {
		o(&opts)
	}

	strip, err := pixarray.NewWS281x(
		opts.NumPixels,
		3,
		opts.Order,
		uint(opts.OscFrequency),
		opts.DMAChannel,
		opts.PWMPins,
	)
	if err != nil {
		return nil, err
	}

	if opts.Brightness != 0. {
		return &LEDStripWithBrightness{
			LEDStrip:   strip,
			Brightness: opts.Brightness,
		}, nil
	} else {
		return strip, nil
	}
}
