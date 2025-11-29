package main

import (
	"flag"
	"fmt"
	"net/http"
)

type Server struct {}

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

	sh, err := NewStatusHandler(
		*numPixels,
		*order,
		float32(*brightness),
		*crossfadeMs,
	)
	if err != nil {
		panic(err)
	}
	sh.InitializeServer(nil)
	http.ListenAndServe(fmt.Sprintf("%s:%d", *listenHost, *listenPort), nil)
}
