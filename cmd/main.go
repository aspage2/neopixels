package main

import (
	"encoding/json"
	"fmt"
	"leds/internal"
	"net/http"

	"github.com/Jon-Bright/ledctl/pixarray"
)

const Brightness = 0.3

func color(r, g, b int) pixarray.Pixel {
	return pixarray.Pixel{
		R: r,
		G: g,
		B: b,
		W: 0,
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
		s.arr.SetAll(Scale(color(c[0], c[1], c[2]), Brightness))
	} else {
		rw.WriteHeader(400)
		fmt.Fprint(rw, "must define 'color'")
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
	pixels, err := internal.NewStrip(39, internal.WithPixelOrder(pixarray.RGB))
	if err != nil {
		panic(err)
	}

	arr := pixarray.NewPixArray(39, 3, pixels)

	arr.SetAll(Scale(color(0, 0, 0), Brightness))
	arr.Write()

	s := Server{arr: arr}

	http.HandleFunc("/set/", s.Set)
	http.HandleFunc("/off/", s.Off)

	http.ListenAndServe("127.0.0.1:5000", nil)
}
