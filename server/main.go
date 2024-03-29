package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Jon-Bright/ledctl/pixarray"
)

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


type Server struct {
	arr *pixarray.PixArray

	order string

	artistMu   *sync.RWMutex
	artistType ArtistType
	artist     Artist

	crossfade int
}


func (s *Server) Status(rw http.ResponseWriter, req *http.Request) {
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

		s.artistType = typ
		s.artist = a
		newPattern := make([]pixarray.Pixel, s.arr.NumPixels())
		a.Draw(newPattern)

		if s.crossfade > 0 {
			old := s.arr.GetPixels()

			const waitMs = 3

			for t := 0; t < s.crossfade; t += waitMs {
				for i := 0; i < len(old); i ++ {
					s.arr.SetOne(i, lerp(old[i], newPattern[i], float32(t) / float32(s.crossfade)))
				}
				s.arr.Write()
				time.Sleep(waitMs * time.Millisecond)
			}
		} else {
			for i := 0; i < len(newPattern); i ++ {
				s.arr.SetOne(i, newPattern[i])
			}
		}

	}

	json.NewEncoder(rw).Encode(envelope)
}

func (s *Server) Off(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	rw.WriteHeader(201)
	s.arr.SetAll(pixarray.Pixel{})
	s.arr.Write()
}

func (s *Server) GetServerSettings(rw http.ResponseWriter, req *http.Request) {
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

func AllowMethods(h http.Handler, method ...string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		found := false
		for _, m := range method {
			if req.Method == m {
				found = true
				break
			}
		}
		if !found {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(rw, req)
	})
}

func main() {

	numPixels := flag.Int("n", 30, "number of leds in the strip")
	listenHost := flag.String("host", "127.0.0.1", "host ip to listen on")
	listenPort := flag.Int("port", 4000, "port to listen on")
	crossfadeMs := flag.Int("xfade", 300, "crossfade time between patterns in milliseconds")
	order := flag.String("order", "grb", "color order for leds [rgb, grb]")
	brightness := flag.Float64("brightness", 0.7, "max brightness of the LEDS")

	flag.Parse()

	pixels, err := NewStrip(
		*numPixels,
		WithPixelOrder(*order),
		WithBrightness(float32(*brightness)),
	)
	if err != nil {
		panic(err)
	}

	arr := pixarray.NewPixArray(*numPixels, 3, pixels)

	s := Server{
		arr: arr,

		artistMu:   &sync.RWMutex{},
		artistType: SOLID,
		artist:     &Solid{Color: [3]int{0, 0, 0}},

		order: *order,

		crossfade: *crossfadeMs,
	}

	http.Handle("/status/", AllowMethods(
		http.HandlerFunc(s.Status),
		http.MethodPost,
		http.MethodGet,
	))

	http.Handle("/settings/", AllowMethods(
		http.HandlerFunc(s.GetServerSettings), http.MethodGet,
	))

	http.Handle("/off/", AllowMethods(
		http.HandlerFunc(s.Off),
		http.MethodPost,
	))

	http.ListenAndServe(fmt.Sprintf("%s:%d", *listenHost, *listenPort), nil)
}
