package main

import (
	"github.com/Jon-Bright/ledctl/pixarray"
)

type Pattern interface {
	Update([]pixarray.Pixel)
}

type PatternFunc func([]pixarray.Pixel)

func (pf PatternFunc) Update(buffer []pixarray.Pixel) {
	pf(buffer)
}

func Solid(px pixarray.Pixel) Pattern {
	return PatternFunc(func(buffer []pixarray.Pixel) {
		for i := 0; i < len(buffer); i++ {
			buffer[i] = px
		}
	})
}
