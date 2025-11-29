//go:build arm

package main

import (
	"fmt"
	"strings"
	"encoding/json"
	"net/http"
	"sync"
	"time"

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
	for i := 0; i < len(arr); i++ {
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

func color(c []int) pixarray.Pixel {
	return pixarray.Pixel{
		R: c[0],
		G: c[1],
		B: c[2],
		W: 0,
	}
}

func lerp(c1, c2 pixarray.Pixel, t float32) pixarray.Pixel {
	return pixarray.Pixel{
		R: c1.R + int(t*float32(c2.R-c1.R)),
		G: c1.G + int(t*float32(c2.G-c1.G)),
		B: c1.B + int(t*float32(c2.B-c1.B)),
		W: c1.W + int(t*float32(c2.W-c1.W)),
	}
}

func Scale(c1 pixarray.Pixel, t float32) pixarray.Pixel {
	return pixarray.Pixel{
		R: int(t * float32(c1.R)),
		G: int(t * float32(c1.G)),
		B: int(t * float32(c1.B)),
		W: int(t * float32(c1.W)),
	}
}

type StatusHandler struct {
	arr        *pixarray.PixArray
	artistMu   *sync.RWMutex
	artistType ArtistType
	artist     Artist
	order      string
	crossfade  int
}

func (s *StatusHandler) Status(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var envelope struct {
		Type ArtistType      `json:"type"`
		Data json.RawMessage `json:"data"`
	}

	switch req.Method {

	case http.MethodGet:
		s.artistMu.RLock()
		defer s.artistMu.RUnlock()

		envelope.Type = s.artistType
		d, _ := json.Marshal(s.artist)
		envelope.Data = d

	case http.MethodPost:
		s.artistMu.Lock()
		defer s.artistMu.Unlock()
		err := json.NewDecoder(req.Body).Decode(&envelope)
		if err != nil {
			rw.WriteHeader(400)
			fmt.Fprintf(rw, "can't parse json body: %s", err)
			return
		}
		var (
			typ ArtistType = envelope.Type.Normalize()
			a   Artist
		)
		switch typ {
		case SOLID:
			a = new(Solid)
		case SEQUENCE:
			a = new(Sequence)
		case GRADIENT:
			a = new(Gradient)
		default:
			rw.WriteHeader(400)
			fmt.Fprintf(rw, "invalid pattern type: %s", envelope.Type)
			return
		}

		err = json.Unmarshal(envelope.Data, a)
		if err != nil {
			rw.WriteHeader(400)
			fmt.Fprintf(rw, "can't parse json body: %s", err)
			return
		}
		
		old := make([]pixarray.Pixel, s.arr.NumPixels())
		s.artist.Draw(old)

		s.artistType = typ
		s.artist = a
		newPattern := make([]pixarray.Pixel, s.arr.NumPixels())
		a.Draw(newPattern)

		if s.crossfade > 0 {
			const waitMs = 3

			for t := 0; t < s.crossfade; t += waitMs {
				for i := 0; i < len(old); i++ {
					s.arr.SetOne(i, lerp(old[i], newPattern[i], float32(t)/float32(s.crossfade)))
				}
				s.arr.Write()
				time.Sleep(waitMs * time.Millisecond)
			}
		} else {
			for i := 0; i < len(newPattern); i++ {
				s.arr.SetOne(i, newPattern[i])
			}
		}

	}

	json.NewEncoder(rw).Encode(envelope)
}

func (s *StatusHandler) Off(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	rw.WriteHeader(201)
	s.arr.SetAll(pixarray.Pixel{})
	s.arr.Write()
}

func (s *StatusHandler) GetServerSettings(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var resp struct {
		NumLeds  int    `json:"numLeds"`
		Channels string `json:"channels"`
		Order    string `json:"order"`
	}
	resp.NumLeds = s.arr.NumPixels()
	if s.arr.NumColors() == 3 {
		resp.Channels = s.order
	} else {
		resp.Channels = s.order + "A"
	}
	resp.Order = s.order
	json.NewEncoder(rw).Encode(resp)
}

func (s *StatusHandler) InitializeServer(sm *http.ServeMux) {
	if sm == nil {
		sm = http.DefaultServeMux
	}

	http.Handle("/lights/settings", AllowMethods(
		http.HandlerFunc(s.GetServerSettings), http.MethodGet,
	))

	http.Handle("/lights/off", AllowMethods(
		http.HandlerFunc(s.Off),
		http.MethodPost,
	))

	http.Handle("/lights/status", AllowMethods(
		http.HandlerFunc(s.Status),
		http.MethodPost,
		http.MethodGet,
	))
}

func NewStatusHandler(
	numPixels int,
	pixelOrder string,
	brightness float32,
	xfade int,
) (*StatusHandler, error) {
	pixels, err := NewStrip(
		numPixels,
		WithPixelOrder(pixelOrder),
		WithBrightness(float32(brightness)),
	)
	arr := pixarray.NewPixArray(numPixels, 3, pixels)
	if err != nil {
		return nil, err
	}
	arr.SetAll(pixarray.Pixel{R: 0, G: 0, B: 0})
	arr.Write()
	return &StatusHandler{
		arr: arr,

		artistMu:   &sync.RWMutex{},
		artistType: SOLID,
		artist:     &Solid{Color: [3]int{0, 0, 0}},

		order: pixelOrder,
		crossfade: xfade,
	}, nil
}
