package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"time"
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
	timestampsFilePath := filepath.Join(directory, "timestamps.txt")

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

	timestampsFile, err := os.Create(timestampsFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer timestampsFile.Close()

	// Variables pour le calcul des timestamps
	var cumulativeDuration time.Duration
	transitionDuration := 2 * time.Second

	transitionImagePath, transitionVideoPath := transition.CreateTransitionVideo(directory)

	for _, video := range videos {
		// Écrivez le timestamp dans "timestamps.txt"
		humanReadableTime := formatDuration(cumulativeDuration)
		timestampsFile.WriteString(fmt.Sprintf("%s - %s\n", humanReadableTime, filepath.Base(video)))

		// Obtenez la durée de la vidéo
		duration, err := getVideoDuration(video)
		if err != nil {
			log.Fatalf("Failed to get duration for %s: %v", video, err)
		}

		// Ajoutez la durée de la vidéo à la durée cumulée
		cumulativeDuration += duration

		file.WriteString(fmt.Sprintf("file '%s'\n", video))
		file.WriteString(fmt.Sprintf("file '%s'\n", transitionVideoPath))
		cumulativeDuration += transitionDuration
	}

	fmt.Println("Merging videos...")

	videoInputOpt := ffmpeg.KwArgs{"f": "concat", "safe": 0}
	// videoOutputOpt := ffmpeg.KwArgs{"c": "copy"}
	// videoOutputOpt := ffmpeg.KwArgs{"preset": "fast", "c": "copy", "pix_fmt": "yuv420p"}
	videoOutputOpt := ffmpeg.KwArgs{"preset": "fast", "c:v": "libx264", "c:a": "aac", "crf": 24, "pix_fmt": "yuv420p", "movflags": "faststart"}
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

func getVideoDuration(videoPath string) (time.Duration, error) {
	probe, err := ffmpeg.Probe(videoPath)
	if err != nil {
		return 0, err
	}

	// Extraire la durée depuis les métadonnées
	durationSeconds, err := probeDuration(probe)

	// durationStr := probe.Get("format").Get("duration").String()
	// durationSeconds, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}

	return time.Duration(durationSeconds * float64(time.Second)), nil
}

type probeFormat struct {
	Duration string `json:"duration"`
}

type probeData struct {
	Format probeFormat `json:"format"`
}

func probeDuration(a string) (float64, error) {
	pd := probeData{}
	err := json.Unmarshal([]byte(a), &pd)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(pd.Format.Duration, 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

// Formate une durée en HH:MM:SS
func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
