package main

import (
	"crypto/sha256"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"time"

	"github.com/CimimUxMaio/artscii"
	"github.com/akamensky/argparse"
	"github.com/eiannone/keyboard"
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

	event, err := keyboard.GetKeys(1)
	checkError(err)
	defer func() { _ = keyboard.Close() }()

	end := false
	for !end {
		webcam.Read(&camImg)
		gocv.Resize(camImg, &camImg, image.Point{*width, *height}, 0, 0, gocv.InterpolationNearestNeighbor)
		img, err := camImg.ToImage()
		checkError(err)
		ascii := artscii.FromImage(img, []rune(*asciiScale))
		//ascii.Print()

		select {
		case keyEvent := <-event:
			switch keyEvent.Key {
			case keyboard.KeySpace:
				name := generatePhotoName()
				fmt.Println(name)
				_, err = ascii.ToFile(name)
				checkError(err)
			case keyboard.KeyCtrlC:
				end = true
			}
		default:
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func generatePhotoName() string {
	workingDir, err := os.Getwd()
	checkError(err)
	hash := sha256.New()
	hash.Write([]byte(time.Now().String()))
	digest := fmt.Sprintf("%x", hash.Sum(nil))
	fmt.Println(string(digest[:8]))
	return workingDir + "/photo_" + string(digest[:8])
}
