package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"video_merger/config"
	"video_merger/timecodes"
	"video_merger/transition"
)

func main() {
	config.Init()
	mergeWithTimeCodes()
	config.Cleanup()
	fmt.Println("Press 'Enter' to quit")
	var input string
	fmt.Scanln(&input)
}

func mergeWithTimeCodes() {
	VALID_VIDEO_FORMATS := []string{".mkv", ".mp4", ".mov"}

	videoResultPath := filepath.Join(config.CURRENT_DIRECTORY, "result.mkv")
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
		// fmt.Println(videoWithPath)
		videos = append(videos, videoWithPath)
	}

	if len(videos) == 0 {
		log.Fatalln("No video found")
	}

	fmt.Printf("%d videos found.\n", len(videos))

	// <-- Merging videos
	err = merge(videos, videoResultPath)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Video available at: %s\n", videoResultPath)
	// ->

	// <-- Generating timecodes.txt
	timeCodesFile, err := os.Create(timeCodeFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer timeCodesFile.Close()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Generating timecodes...")
	err = timecodes.GenerateTimeCodes(timeCodeFilePath, videos)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Timecodes file available at: %s\n", timeCodeFilePath)

	// ->

	fmt.Println("Done")

}

func merge(videos []string, videoResultPath string) error {
	videoTextFilePath := filepath.Join(config.TEMP_DIRECTORY, "videos.txt")

	file, err := os.Create(videoTextFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	transitionImagePath, transitionVideoPath := transition.CreateTransitionVideo()

	// Calculer duree totale pour le pourcentage
	totalDuration := 0.0
	for _, video := range videos {
		d, _ := timecodes.GetVideoDuration(video)
		totalDuration += d.Seconds()
	}
	// Ajouter les transitions (sauf apres la derniere)
	totalDuration += float64(len(videos)-1) * config.TRANSITION_DURATION.Seconds()

	for i, video := range videos {
		file.WriteString(fmt.Sprintf("file '%s'\n", video))
		if i < len(videos)-1 {
			file.WriteString(fmt.Sprintf("file '%s'\n", transitionVideoPath))
		}
	}
	file.Close()

	fmt.Println("Merging videos...")

	// Construire la commande FFmpeg
	args := []string{
		"-f", "concat", "-safe", "0", "-i", videoTextFilePath,
		"-preset", "fast", "-c:v", "libx264", "-c:a", "aac",
		"-crf", "24", "-pix_fmt", "yuv420p", "-movflags", "faststart",
		"-y", videoResultPath,
	}

	cmd := exec.Command("ffmpeg", args...)

	// Pipe pour capturer stderr
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Parser la progression
	scanner := bufio.NewScanner(stderrPipe)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "time=") {
			parts := strings.Split(line, "time=")
			if len(parts) > 1 {
				timeStr := strings.Split(parts[1], " ")[0]
				currentSeconds := parseTimeToSeconds(timeStr)
				percent := int((currentSeconds / totalDuration) * 100)
				if percent > 100 {
					percent = 100
				}
				fmt.Printf("\rMerging videos... %d%%", percent)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}
	fmt.Println()

	// Nettoyage
	os.Remove(videoTextFilePath)
	os.Remove(transitionImagePath)
	os.Remove(transitionVideoPath)

	return nil
}

// Helper pour parser "HH:MM:SS.xx" en secondes
func parseTimeToSeconds(timeStr string) float64 {
	parts := strings.Split(timeStr, ":")
	var seconds float64
	if len(parts) == 3 {
		h, _ := strconv.ParseFloat(parts[0], 64)
		m, _ := strconv.ParseFloat(parts[1], 64)
		s, _ := strconv.ParseFloat(parts[2], 64)
		seconds = h*3600 + m*60 + s
	}
	return seconds
}
