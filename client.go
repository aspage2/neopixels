package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
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

func colorFromString(arg string) [3]int64 {
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

	return [3]int64{
		(val >> 16) & 0xff,
		(val >> 8) & 0xff,
		val & 0xff,
	}
}

func colorFrom3Parts(args []string) [3]int64 {
	return [3]int64{
		Must(strconv.ParseInt(args[0], 10, 64)).(int64),
		Must(strconv.ParseInt(args[1], 10, 64)).(int64),
		Must(strconv.ParseInt(args[2], 10, 64)).(int64),
	}
}

func swap(arr []int64, i, j int) {
	tmp := arr[i]
	arr[i] = arr[j]
	arr[j] = tmp
}

func main() {
	off := flag.Bool("off", false, "Turn the leds off")
	mode := flag.String("mode", "grb", "either 'grb' or 'rgb'")
	flag.Parse()

	if *off {
		http.Post("http://192.168.2.6:5000/off/", "application/json", bytes.NewReader([]byte("{}")))
		return
	}

	args := flag.Args()

	var color [3]int64
	if len(args) == 1 {
		color = colorFromString(args[0])
	} else if len(args) == 3 {
		color = colorFrom3Parts(args)
	} else {
		panic(errors.New("must provide either color code or 3 args"))
	}

	if *mode == "grb" {
		swap(color[:], 0, 1)
	}

	data, _ := json.Marshal(map[string][3]int64{"color": color})

	resp, err := http.Post("http://192.168.2.6:5000/set/", "application/json", bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
