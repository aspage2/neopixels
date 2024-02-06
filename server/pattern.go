package main

import (
	"strings"

	"github.com/Jon-Bright/ledctl/pixarray"
)

const (
	SOLID    ArtistType = "SOLID"
	SEQUENCE ArtistType = "SEQUENCE"
	GRADIENT ArtistType = "GRADIENT"
)

type Artist interface {
	Draw([]pixarray.Pixel)
}

type ArtistType string

func (at ArtistType) Normalize() ArtistType {
	return ArtistType(strings.ToUpper(string(at)))
}


type Solid struct {
	Color [3]int `json:"color"`
}

func (s *Solid) Draw(arr []pixarray.Pixel) {
	for i := 0; i < len(arr); i ++ {
		arr[i] = color(s.Color[:])
	}
}

type Gradient struct {
	Colors [][3]int `json:"colors"`
}

func (g *Gradient) Draw(arr []pixarray.Pixel) {
	l := len(arr)
	c1 := color(g.Colors[0][:])
	c2 := color(g.Colors[1][:])
	for i := 0; i < l; i++ {
		t := float32(i) / float32(l)
		arr[i] = lerp(c1, c2, t)
	}

}

type Sequence struct {
	Colors [][3]int `json:"colors"`
}

func (seq *Sequence) Draw(arr []pixarray.Pixel) {
	for i := 0; i < len(arr); i++ {
		c := seq.Colors[i%len(seq.Colors)][:]
		arr[i] = color(c)
	}
}
