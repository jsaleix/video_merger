package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	TRANSITION_DURATION = 2 * time.Second
)

var CURRENT_DIRECTORY string
var TEMP_DIRECTORY string

func Init() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error : %v", err)
	}

	absPath, err := filepath.Abs(currentDir)
	if err != nil {
		log.Fatalf("Error with filepath.Abs : %v", err)
	}

	temp_dir, err := os.MkdirTemp("./", "temp")

	if err != nil {
		log.Fatalf("Erreur lors de la conversion en chemin absolu : %v", err)
	}

	CURRENT_DIRECTORY = absPath
	TEMP_DIRECTORY = filepath.Join(CURRENT_DIRECTORY, temp_dir)

	if runtime.GOOS == "windows" {
		CURRENT_DIRECTORY = filepath.ToSlash(CURRENT_DIRECTORY)
		TEMP_DIRECTORY = filepath.ToSlash(TEMP_DIRECTORY)
	}

}

func Cleanup() {
	os.Remove(TEMP_DIRECTORY)
}
