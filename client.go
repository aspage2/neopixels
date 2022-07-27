package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func Must(v interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return v
}

type Color [3]int64

func (c Color) GRB() Color {
	return Color{
		c[1],
		c[0],
		c[2],
	}
}

func colorFromString(arg string) Color {
	var payload string
	if strings.HasPrefix(arg, "0x") {
		if len(arg) != 8 {
			panic("bad arg")
		}
		payload = arg[2:]
	} else if strings.HasPrefix(arg, "#") {
		if len(arg) != 7 {
			panic("bad arg")
		}
		payload = arg[1:]
	}

	val, err := strconv.ParseInt(payload, 16, 64)
	if err != nil {
		panic(err)
	}

	return Color{
		(val >> 16) & 0xff,
		(val >> 8) & 0xff,
		val & 0xff,
	}
}

func Do(payload interface{}) {
	data := Must(json.Marshal(payload)).([]byte)
	resp, err := http.Post("http://192.168.2.16/set/", "application/json", bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	resp.Body.Close()
}

func main() {
	off := flag.Bool("off", false, "Turn the leds off")
	mode := flag.String("mode", "solid", "Lighting mode.")
	flag.Parse()

	if *off {
		http.Post("http://192.168.2.16/off/", "application/json", bytes.NewReader([]byte("{}")))
		return
	}

	args := flag.Args()

	var colors []Color
	for _, a := range args {
		colors = append(colors, colorFromString(a))
	}

	switch *mode {
	case "solid":
		if len(colors) > 1 {
			fmt.Printf("WARN: mode is 'solid', only using %s\n", args[0])
		}
		Do(map[string]Color{"solid": colors[0]})
	case "gradient":
		if len(colors) <= 1 {
			panic(errors.New("need 2 colors for gradient"))
		} else if len(colors) > 2 {
			fmt.Printf("WARN: mode is 'gradient', only using %s, %s\n", args[0], args[1])
		}
		Do(map[string][]Color{"gradient": colors[:2]})
	case "sequence":
		if len(colors) == 1 {
			Do(map[string]Color{"solid": colors[0]})
		} else {
			Do(map[string][]Color{"sequence": colors})
		}
	case "merry":
		R := Color{0xff, 0, 0}
		G := Color{0, 0xff, 0}
		B := Color{0, 0, 0xff}
		V := Color{0xff, 0, 0xff}
		Y := Color{0xff, 0xff, 0}
		seq := []Color{G, G, R, G, G, B, G, G, Y, G, G, V}
		Do(map[string][]Color{"sequence": seq})
	}

}
