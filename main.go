package main

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/CimimUxMaio/artscii"
	"github.com/akamensky/argparse"
	"gocv.io/x/gocv"
)

func main() {
	parser := argparse.NewParser("AsciiCamera", "Shows camera recorded video as ascii characters in command line")
	width := parser.Int("x", "width", &argparse.Options{Required: true, Help: "Image width"})
	height := parser.Int("y", "height", &argparse.Options{Required: true, Help: "Image height"})
	defaultScale := " `-:~*r+=xhwAD9MWB@"
	asciiScale := parser.String("s", "scale", &argparse.Options{Required: false, Default: defaultScale, Help: "Ascii characters to use ordered by brightness (from low to high)"})
	err := parser.Parse(os.Args)
	checkError(err)

	webcam, err := gocv.VideoCaptureDevice(0)
	checkError(err)
	defer webcam.Close()

	camImg := gocv.NewMat()
	defer camImg.Close()

	for {
		webcam.Read(&camImg)
		gocv.Resize(camImg, &camImg, image.Point{*width, *height}, 0, 0, gocv.InterpolationNearestNeighbor)
		img, err := camImg.ToImage()
		checkError(err)
		ascii := artscii.FromImage(img, []byte(*asciiScale))
		ascii.Print()
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
