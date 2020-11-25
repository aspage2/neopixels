package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"pnpleds/ledproto"
	"pnpleds/ledserver/internal"
	"pnpleds/ledserver/internal/pattern"
	"sync"
	"syscall"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

type LedStripServicer struct {
	PatternC chan pattern.Pattern
	Device   *ws2811.WS2811
}

func NewLedStripServicer(patternC chan pattern.Pattern, dev *ws2811.WS2811) *grpc.Server {
	servicer := &LedStripServicer{
		PatternC: patternC,
		Device:   dev,
	}
	server := grpc.NewServer()
	ledproto.RegisterLedStripService(server, &ledproto.LedStripService{
		Set:         servicer.Set,
		SetRealtime: servicer.SetStream,
	})
	return server
}

func (serv *LedStripServicer) SetStream(stream ledproto.LedStrip_SetRealtimeServer) error {
	leds := serv.Device.Leds(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		internal.Fill(leds, req.Color)
		serv.Device.Render()
	}
}

func (serv *LedStripServicer) Set(
	ctx context.Context, req *ledproto.SetPatternRequest,
) (*empty.Empty, error) {

	return &empty.Empty{}, nil

	//var p pattern.Pattern
	//switch x := req.Pattern.(type) {
	//case *ledproto.SetPatternRequest_Solid:
	//	p = &pattern.Solid{Color: x.Solid.Color}
	//case *ledproto.SetPatternRequest_Stripe:
	//	colors := x.Stripe.Colors
	//	if len(colors) == 0 {
	//		return nil, errors.New("zero colors")
	//	} else if len(colors) == 1 {
	//		p = &pattern.Solid{Color: colors[0]}
	//	} else {
	//		p = &pattern.Stripe{
	//			StripeSize: int(x.Stripe.StripeLen),
	//			Colors:     colors,
	//		}
	//	}
	//case *ledproto.SetPatternRequest_Snowflake:
	//	p = &pattern.Snowflake{ChanceSpawn: float64(x.Snowflake.SpawnChance)}
	//case *ledproto.SetPatternRequest_Breathe:
	//	colors := x.Breathe.Colors
	//	if len(colors) == 0 {
	//		return nil, errors.New("zero colors")
	//	} else if len(colors) == 1 {
	//		p = &pattern.Solid{Color: colors[0]}
	//	} else {
	//		p = &pattern.Breathe{
	//			Colors: colors,
	//		}
	//	}
	//}
	//select {
	//case serv.PatternC <- p:
	//	return &empty.Empty{}, nil
	//case <-ctx.Done():
	//	return nil, nil
	//}
}

func LerpPatterns(frameC, pat1, pat2 chan []uint32, n int) {
	var f1, f2 []uint32
	leds := make([]uint32, n)
	for t := 0.; t < 1.0; t += 0.01 {
		f1 = <-pat1
		f2 = <-pat2
		for i := range leds {
			leds[i] = internal.Lerp(f1[i], f2[i], t)
		}
		frameC <- leds
	}
}

func render(frameC chan []uint32, dev *ws2811.WS2811) {
	leds := dev.Leds(0)
	for frame := range frameC {
		copy(leds, frame)
		dev.Render()
		time.Sleep(20 * time.Millisecond)

	}
}

func ledRenderLoop(patterns chan pattern.Pattern) {
	N := viper.GetInt("NumLeds")
	dev, err := internal.MakeLedDevice(
		N, viper.GetInt("Brightness"),
	)
	chk(err)
	chk(dev.Init())
	defer internal.DeviceCleanup(dev)

	var (
		cancel context.CancelFunc
		ctx    context.Context
		wg     sync.WaitGroup
	)

	frameC := make(chan []uint32)
	defer close(frameC)
	go render(frameC, dev)

	for p := range patterns {
		if cancel != nil {
			cancel()
			wg.Wait()
		}
		ctx, cancel = context.WithCancel(context.Background())
		wg.Add(1)
		go func(pat pattern.Pattern) {
			pat.Pattern(ctx, frameC, N)
			wg.Done()
		}(p)
	}
	if cancel != nil {
		cancel()
	}
}

func run(serv *grpc.Server) {
	socket, err := net.Listen("tcp", "0.0.0.0:50051")
	chk(err)
	serv.Serve(socket)
}

func main() {
	patternC := make(chan pattern.Pattern)
	N := viper.GetInt("NumLeds")
	dev, err := internal.MakeLedDevice(
		N, viper.GetInt("Brightness"),
	)
	chk(err)
	chk(dev.Init())
	defer internal.DeviceCleanup(dev)

	server := NewLedStripServicer(patternC, dev)

	sigC := make(chan os.Signal)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigC
		fmt.Println("Interrupt...")
		server.Stop()
	}()

	//var renderWg sync.WaitGroup
	//renderWg.Add(1)
	//go func() {
	//	ledRenderLoop(patternC)
	//	renderWg.Done()
	//}()
	//patternC <- &pattern.Solid{
	//	Color: 0xfedb00, // experts in wtf you want
	//}
	fmt.Println("start serving")
	run(server)
	close(patternC)
	//renderWg.Wait()
}
