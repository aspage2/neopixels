package internal

import (
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
	"github.com/spf13/viper"
)

// CheckError panics on error, noops otherwise.
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// Clear the full strip of leds. Safe on nil slices
func Clear(leds []uint32) {
	Fill(leds, 0)
}

func Fill(leds []uint32, color uint32) {
	for i := range leds {
		leds[i] = color
	}
}

func DeviceFromConfig() (*ws2811.WS2811, error) {
	return MakeLedDevice(
		viper.GetInt("NumLeds"),
		viper.GetInt("Brightness"),
	)
}

// MakeLedDevice initializes and returns a new WS2811
// device interface for controlling an LED strip
func MakeLedDevice(ledCounts, brightness int) (*ws2811.WS2811, error) {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = brightness
	opt.Channels[0].LedCount = ledCounts

	return ws2811.MakeWS2811(&opt)
}

// DeviceCleanup zeroes-out all leds on a strip and calls
// the device interface's cleanup method.
func DeviceCleanup(dev *ws2811.WS2811) {
	Clear(dev.Leds(0))
	dev.Render()
	dev.Fini()
}
