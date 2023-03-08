package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

const (
	HostBase = "http://192.168.0.232:9000/"
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

func colorFromString(arg string) (Color, error) {
	var payload string
	payload = strings.TrimPrefix(arg, "0x")
	payload = strings.TrimPrefix(payload, "#")
	if len(payload) != 6 {
		return Color{}, errors.New("hex colors must be 6 characters")
	}

	val, err := strconv.ParseInt(payload, 16, 64)
	if err != nil {
		panic(err)
	}

	return Color{
		(val >> 16) & 0xff,
		(val >> 8) & 0xff,
		val & 0xff,
	}, nil
}

func getUrl(path string) string {
	path = strings.TrimLeft(path, "/")
	return strings.TrimRight(HostBase, "/") + "/" + path
}

func Post(path string, data []byte) error {
	resp, err := http.Post(getUrl(path), "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func Get(path string) ([]byte, error) {
	resp, err := http.Get(getUrl(path))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func MakePayload(patternType string, data interface{}) ([]byte, error) {
	var envelope struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	envelope.Type = patternType
	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	envelope.Data = d
	return json.Marshal(envelope)
}

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "off",
				Usage: "turn the leds off",
				Action: func(cliCtx *cli.Context) error {
					return Post("off/", []byte{})
				},
			},
			{
				Name:      "gradient",
				Usage:     "make a gradient pattern",
				ArgsUsage: "colors to include",
				Action: func(ctx *cli.Context) error {
					args := ctx.Args()
					c1, c1Err := colorFromString(args.Get(0))
					if c1Err != nil {
						return c1Err
					}
					c2, c2Err := colorFromString(args.Get(1))
					if c2Err != nil {
						return c2Err
					}
					envelope := struct {
						Colors [2]Color `json:"colors"`
					}{Colors: [2]Color{c1, c2}}
					payload := Must(MakePayload("GRADIENT", envelope)).([]byte)
					return Post("/status/", payload)
				},
			},
			{
				Name:      "sequence",
				Usage:     "make a sequence pattern",
				ArgsUsage: "colors to include",
				Action: func(ctx *cli.Context) error {
					args := ctx.Args()
					c1, c1Err := colorFromString(args.Get(0))
					if c1Err != nil {
						return c1Err
					}
					c2, c2Err := colorFromString(args.Get(1))
					if c2Err != nil {
						return c2Err
					}
					envelope := struct {
						Colors [2]Color `json:"colors"`
					}{Colors: [2]Color{c1, c2}}
					payload := Must(MakePayload("SEQUENCE", envelope)).([]byte)
					return Post("/status/", payload)
				},
			},
			{
				Name:      "solid",
				Usage:     "make a solid pattern",
				ArgsUsage: "colors to include",
				Action: func(ctx *cli.Context) error {
					args := ctx.Args()
					c, err := colorFromString(args.Get(0))
					if err != nil {
						return err
					}
					envelope := struct {
						Color `json:"color"`
					}{Color: c}
					payload := Must(MakePayload("SOLID", envelope)).([]byte)
					return Post("/status/", payload)
				},
			},
			{
				Name:  "get",
				Usage: "get current status",
				Action: func(ctx *cli.Context) error {
					resp, err := Get("/status/")
					if err != nil {
						return err
					}
					fmt.Println(string(resp))
					return nil
				},
			},
			{
				Name:  "merry",
				Usage: "merry xmas",
				Action: func(ctx *cli.Context) error {
					R := Color{0xff, 0, 0}
					G := Color{0, 0xff, 0}
					B := Color{0, 0, 0xff}
					V := Color{0xff, 0, 0xff}
					Y := Color{0xff, 0xff, 0}
					seq := []Color{G, G, R, G, G, B, G, G, Y, G, G, V}

					return Post("/status/", Must(MakePayload("SEQUENCE", map[string][]Color{"colors": seq})).([]byte))
				},
			},
			{
				Name:  "get-settings",
				Usage: "get current settings",
				Action: func(ctx *cli.Context) error {
					resp, err := Get("/settings/")
					if err != nil {
						return err
					}
					fmt.Println(string(resp))
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
