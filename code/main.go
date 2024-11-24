package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"video_merger/transition"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {
	merge()
	fmt.Println("Press 'Enter' to close")
	var input string
	fmt.Scanln(&input)

}

func merge() {
	VALID_VIDEO_FORMATS := []string{".mkv", ".mp4", ".mov"}

	directory := "./"
	outputPath := filepath.Join(directory, "result.mkv")
	videoTextFilePath := filepath.Join(directory, "videos.txt")

	// Listing files
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatalln(err)
		return
	}

	videos := []string{}

	for _, file := range files {
		fileExt := filepath.Ext(file.Name())
		if file.IsDir() || !slices.Contains(VALID_VIDEO_FORMATS, fileExt) || file.Name() == "result.mp4" {
			continue
		}
		videoWithPath := filepath.Join(directory, file.Name())
		fmt.Println(videoWithPath)
		videos = append(videos, videoWithPath)
	}

	if len(videos) == 0 {
		log.Fatalln("No videos found")
	}

	file, err := os.Create(videoTextFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	transitionImagePath, transitionVideoPath := transition.CreateTransitionVideo(directory)

	for _, video := range videos {
		file.WriteString(fmt.Sprintf("file '%s'\n", video))
		file.WriteString(fmt.Sprintf("file '%s'\n", transitionVideoPath))
	}

	fmt.Println("Merging videos...")

	videoInputOpt := ffmpeg.KwArgs{"f": "concat", "safe": 0}
	// videoOutputOpt := ffmpeg.KwArgs{"c": "copy"}
	videoOutputOpt := ffmpeg.KwArgs{"preset": "fast", "c:v": "libx264", "c:a": "aac", "crf": 24, "pix_fmt": "yuv420p", "movflags": "faststart"}
	// videoOutputOpt := ffmpeg.KwArgs{"preset": "fast", "c": "copy", "pix_fmt": "yuv420p"}
	err = ffmpeg.Input("videos.txt", videoInputOpt).Output(outputPath, videoOutputOpt).OverWriteOutput().Run()

	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Printf("Video available at: %s\n", outputPath)
	}

	file.Close()
	os.Remove(videoTextFilePath)
	os.Remove(transitionImagePath)
	os.Remove(transitionVideoPath)
}
