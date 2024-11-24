package timecodes

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type probeFormat struct {
	Duration string `json:"duration"`
}

type probeData struct {
	Format probeFormat `json:"format"`
}

func FormatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func GetVideoDuration(videoPath string) (time.Duration, error) {
	probe, err := ffmpeg.Probe(videoPath)
	if err != nil {
		return 0, err
	}

	durationSeconds, err := probeDuration(probe)

	if err != nil {
		return 0, err
	}

	return time.Duration(durationSeconds * float64(time.Second)), nil
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
