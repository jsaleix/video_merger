package timecodes

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"video_merger/config"
)

func GenerateTimeCodes(timeCodeFilePath string, videosPaths []string) error {
	var cumulativeDuration time.Duration

	timeCodesFile, err := os.Create(timeCodeFilePath)
	if err != nil {
		return err
	}
	defer timeCodesFile.Close()

	timeCodesFile.WriteString("Timecodes\n\n")
	for _, videoPath := range videosPaths {
		humanReadableTime := FormatDuration(cumulativeDuration)
		timeCodesFile.WriteString(fmt.Sprintf("%s - %s\n", humanReadableTime, filepath.Base(videoPath)))

		duration, err := GetVideoDuration(videoPath)
		if err != nil {
			return fmt.Errorf("failed to get duration for %s: %v", videoPath, err)
		}

		cumulativeDuration += duration
		cumulativeDuration += config.TRANSITION_DURATION
	}

	return nil

}
