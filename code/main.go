package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"video_merger/config"
	"video_merger/timecodes"
	"video_merger/transition"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {
	config.Init()
	mergeWithTimeCodes()
	fmt.Println("Press 'Enter' to close")
	var input string
	fmt.Scanln(&input)
}

func mergeWithTimeCodes() {
	VALID_VIDEO_FORMATS := []string{".mkv", ".mp4", ".mov"}

	videoResultPath := filepath.Join(config.CURRENT_DIRECTORY, "result.mkv")
	videoTextFilePath := filepath.Join(config.TEMP_DIRECTORY, "videos.txt")
	timeCodeFilePath := filepath.Join(config.CURRENT_DIRECTORY, "timecodes.txt")

	// Listing files
	files, err := os.ReadDir(config.CURRENT_DIRECTORY)
	if err != nil {
		log.Fatalln(err)
		return
	}

	videos := []string{}

	for _, file := range files {
		fileExt := filepath.Ext(file.Name())
		if file.IsDir() || !slices.Contains(VALID_VIDEO_FORMATS, fileExt) || file.Name() == "result.mkv" {
			continue
		}
		videoWithPath := filepath.Join(config.CURRENT_DIRECTORY, file.Name())
		fmt.Println(videoWithPath)
		videos = append(videos, videoWithPath)
	}

	if len(videos) == 0 {
		log.Fatalln("No videos found")
	}

	err = merge(videos, videoTextFilePath, videoResultPath)

	if err != nil {
		log.Println(err)
	}

	timeCodesFile, err := os.Create(timeCodeFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer timeCodesFile.Close()

	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Printf("Video available at: %s\n", videoResultPath)
	}

	fmt.Println("Generating timecodes...")
	err = timecodes.GenerateTimeCodes(timeCodeFilePath, videos)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Timecodes generated successfully")

}

func merge(videos []string, videoTextFilePath string, videoResultPath string) error {

	file, err := os.Create(videoTextFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	transitionImagePath, transitionVideoPath := transition.CreateTransitionVideo()

	for _, video := range videos {
		file.WriteString(fmt.Sprintf("file '%s'\n", video))
		file.WriteString(fmt.Sprintf("file '%s'\n", transitionVideoPath))
	}

	fmt.Println("Merging videos...")

	videoInputOpt := ffmpeg.KwArgs{"f": "concat", "safe": 0}
	videoOutputOpt := ffmpeg.KwArgs{"c": "copy"}
	// videoOutputOpt := ffmpeg.KwArgs{"preset": "fast", "c:v": "libx264", "c:a": "aac", "crf": 24, "pix_fmt": "yuv420p", "movflags": "faststart"}
	err = ffmpeg.Input(videoTextFilePath, videoInputOpt).Output(videoResultPath, videoOutputOpt).OverWriteOutput().Run()

	file.Close()
	os.Remove(videoTextFilePath)
	os.Remove(transitionImagePath)
	os.Remove(transitionVideoPath)

	return err

}
