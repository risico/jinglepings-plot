package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strings"
	"time"

	fastping "github.com/tatsushid/go-fastping"
)

// IP prefix of jinglepings
const ipPrefix = "2001:4c08:2028"
const maxRTT = time.Millisecond * 100
const sleep = 50 * time.Millisecond

func main() {
	// You can register another format here
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	if os.Args[1] == "" {
		fmt.Println("Error: please set an image")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])

	if err != nil {
		fmt.Println("Error: File could not be opened")
		os.Exit(1)
	}

	defer file.Close()

	pixels, err := getPixels(file)

	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}

	pingPixels(pixels)
}

func pingAll(pixels []Pixel) {
	p := fastping.NewPinger()
	p.Network("udp")
	p.MaxRTT = maxRTT

	for x := 0; x < len(pixels); x++ {
		p.AddIP(pixels[x].toIPv6())
	}

	err := p.Run()
	if err != nil {
		fmt.Printf("ERRRO: %+v \n", err)
	}
}

func pingPixels(pixels []Pixel) {
	for {
		go pingAll(pixels)

		time.Sleep(sleep)
	}
}

// Get the bi-dimensional pixel array
func getPixels(file io.Reader) ([]Pixel, error) {

	img, _, err := image.Decode(file)

	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels []Pixel
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := rgbaToPixel(img.At(x, y).RGBA())
			pixel.X = x
			pixel.Y = y

			if pixel.R != 0 && pixel.G != 0 && pixel.B != 0 {
				pixels = append(pixels, pixel)
			}
		}
	}

	return pixels, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{R: int(r / 257), G: int(g / 257), B: int(b / 257)}
}

func (c *Pixel) toIPv6() string {
	return strings.ToUpper(fmt.Sprintf("%s:%d:%d:%02x:%02x:%02x", ipPrefix, c.X, c.Y, c.R, c.G, c.B))
}

// Pixel struct example
type Pixel struct {
	X int
	Y int
	R int
	G int
	B int
}
