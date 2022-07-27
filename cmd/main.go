package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"leds/internal"
	"net/http"

	"github.com/Jon-Bright/ledctl/pixarray"
)

const Brightness = 0.3

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
}

func (s *Server) Set(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var reqEnvelope struct {
		Color    []int
		Sequence [][]int
		Gradient [][]int
	}

	err := json.NewDecoder(req.Body).Decode(&reqEnvelope)
	if err != nil {
		rw.WriteHeader(400)
		fmt.Fprintf(rw, "can't parse json body: %s", err)
		return
	}

	if c := reqEnvelope.Color; c != nil {
		if len(c) != 3 {
			rw.WriteHeader(400)
			fmt.Fprint(rw, "colors must be lists of 3 integers")
			return
		}
		s.arr.SetAll(Scale(color(c), Brightness))
	} else if seq := reqEnvelope.Sequence; seq != nil {
		for _, c := range seq {
			if len(c) != 3 {
				rw.WriteHeader(400)
				fmt.Fprint(rw, "colors must be lists of 3 integers")
				return
			}
		}
		for i := 0; i < s.arr.NumPixels(); i++ {
			c := seq[i%len(seq)]
			s.arr.SetOne(i, Scale(color(c), Brightness))
		}
	} else if grad := reqEnvelope.Gradient; grad != nil {
		for _, c := range seq {
			if len(c) != 3 {
				rw.WriteHeader(400)
				fmt.Fprint(rw, "colors must be lists of 3 integers")
				return
			}
		}
		c1 := color(grad[0])
		c2 := color(grad[1])
		for i := 0; i < s.arr.NumPixels(); i++ {
			t := float32(i) / float32(s.arr.NumPixels())
			s.arr.SetOne(i, Scale(lerp(c1, c2, t), Brightness))
		}

	} else {
		rw.WriteHeader(400)
		fmt.Fprint(rw, "must define 'color', 'sequence' or 'gradient'")
		return
	}
	rw.WriteHeader(201)
	s.arr.Write()
}

func (s *Server) Off(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	rw.WriteHeader(201)
	s.arr.SetAll(pixarray.Pixel{})
	s.arr.Write()
}
func main() {

	numPixels := flag.Int("n", 30, "number of leds in the strip")
	listenHost := flag.String("host", "127.0.0.1", "host ip to listen on")
	listenPort := flag.Int("port", 4000, "port to listen on")
	order := flag.String("order", "grb", "color order for leds [rgb, grb]")

	flag.Parse()

	pixels, err := internal.NewStrip(*numPixels, internal.WithPixelOrder(pixarray.StringOrders[*order]))
	if err != nil {
		panic(err)
	}

	arr := pixarray.NewPixArray(*numPixels, 3, pixels)

	arr.SetAll(Scale(pixarray.Pixel{}, Brightness))
	arr.Write()

	s := Server{arr: arr}

	http.HandleFunc("/set/", s.Set)
	http.HandleFunc("/off/", s.Off)

	http.ListenAndServe(fmt.Sprintf("%s:%d", *listenHost, *listenPort), nil)
}
