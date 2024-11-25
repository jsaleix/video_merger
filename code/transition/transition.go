package transition

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"
	"video_merger/config"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func CreateTransitionVideo() (imgPath, videoPath string) {
	transitionImagePath := filepath.Join(config.TEMP_DIRECTORY, "transition.jpg")
	transitionVideoPath := filepath.Join(config.TEMP_DIRECTORY, "transition.mkv")

	createTransitionImage(transitionImagePath)

	// fmt.Println(transitionVideoPath)
	inputOpts := ffmpeg.KwArgs{"f": "image2", "loop": 1, "framerate": 1, "t": config.TRANSITION_DURATION}
	// inputOpts := ffmpeg.KwArgs{"loop": 1}
	outputOpts := ffmpeg.KwArgs{"c:v": "libx264", "c:a": "aac", "crf": 18, "pix_fmt": "yuv420p", "movflags": "faststart"}
	// outputOpts := ffmpeg.KwArgs{"c": "libx264"}
	ffmpeg.Input(transitionImagePath, inputOpts).Output(transitionVideoPath, outputOpts).OverWriteOutput().Run()

	time.Sleep(2 * time.Second)
	fmt.Println("Transition video created.")
	return transitionImagePath, transitionVideoPath
}

func createTransitionImage(path string) {
	img := generateImage()
	saveImage(img, path)
}

func saveImage(img *image.RGBA, path string) {
	out, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	err = jpeg.Encode(out, img, nil)
	if err != nil {
		panic(err)
	}

	out.Close()
}

func generateImage() *image.RGBA {
	rect := image.Rect(0, 0, 1920, 1080)
	img := image.NewRGBA(rect)
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 1}}, image.Point{}, draw.Src)

	return img
}
