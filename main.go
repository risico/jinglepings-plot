package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"net"
	"os"
	"time"

	fastping "github.com/tatsushid/go-fastping"
)

const ipPrefix = "2001:4c08:2028"

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

	err = pingImage(file)

	if err != nil {
		fmt.Println("Error: Image could not be decoded")
		os.Exit(1)
	}
}

// Get the bi-dimensional pixel array
func pingImage(file io.Reader) error {
	p := fastping.NewPinger()
	p.Network("udp")
	p.MaxRTT = 1

	img, _, err := image.Decode(file)

	if err != nil {
		return err
	}

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		fmt.Printf("|\n")
		for x := 0; x < width; x++ {
			pixel := rgbaToPixel(img.At(x, y).RGBA())
			if pixel.R != 0 && pixel.G != 0 && pixel.B != 0 {
				fmt.Printf("#")

				ip := pixel.toIPv6(x, y)

				ra, err := net.ResolveIPAddr("ip6:icmp", ip)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				p.AddIPAddr(ra)

				row = append(row, pixel)
			} else {
				fmt.Printf(" ")
			}
		}
		pixels = append(pixels, row)
	}

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {}
	p.OnIdle = func() {}
	err = p.Run()
	if err != nil {
	}

	return nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func (c *Pixel) toIPv6(x, y int) string {
	return fmt.Sprintf("%s:%d:%d:%02x:%02x:%02x", ipPrefix, x, y, c.R, c.G, c.B)
}

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}
