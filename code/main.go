package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {
	VALID_VIDEO_FORMATS := []string{".mkv", ".mp4", ".mov"}

	directory := "./"
	outputPath := filepath.Join(directory, "result.mp4")
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

	for _, video := range videos {
		file.WriteString(fmt.Sprintf("file '%s'\n", video))
	}

	err = ffmpeg.Input("videos.txt", ffmpeg.KwArgs{"f": "concat", "safe": 0}).Output(outputPath, ffmpeg.KwArgs{"c": "copy"}).OverWriteOutput().Run()

	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Printf("Video available at: %s\n", outputPath)
	}
	file.Close()

	err = os.Remove(videoTextFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Press 'Enter' to close")
	var input string
	fmt.Scanln(&input)

}
