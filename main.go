package main

import (
	"crypto/sha256"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"time"

	"github.com/CimimUxMaio/artscii/artscii"
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
	defer keyboard.Close()

	end := false
	flash := false
	flashState := 0
	flashDuration := 5
	for !end {
		webcam.Read(&camImg)
		gocv.Resize(camImg, &camImg, image.Point{*width, *height}, 0, 0, gocv.InterpolationNearestNeighbor)

		if flash {
			k := 1 + 20*math.Sin(float64(flashState)*math.Pi/float64(flashDuration))
			camImg.MultiplyFloat(float32(k))
			flashState += 1
			if flashState > flashDuration {
				flashState = 0
				flash = false
			}
		}

		img, err := camImg.ToImage()
		checkError(err)

		ascii := artscii.FromImage(img, []rune(*asciiScale))
		ascii.Print()

		select {
		case keyEvent := <-event:
			switch keyEvent.Key {
			case keyboard.KeySpace:
				flash = true
				_, err = ascii.ToFile(generatePhotoName())
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
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(time.Now().String())))
	return workingDir + "/photo_" + string(hash[:8])
}
