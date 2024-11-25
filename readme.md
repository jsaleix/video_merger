# Video Merger

**Video Merger** is a Go application that merges all video files in the current directory into a single video file, with a 2 seconds gap between each video. Additionally, it generates a `timecode.txt` file that lists the start time of each video along with its name.

---

## Features

- Detects files using `.mp4`, `.mov` and `.mkv` video formats in the current directory.
- Merges all videos into a single output file with a gap of 2 seconds between each video.
- Generates a `timecode.txt` file with the format:
  ```
  Timecodes

  00:00:00 - video_name_1.mp4
  00:05:23 - video_name_2.mov
  ```

---

## Requirements

- [Go](https://go.dev/) installed on your system.
- [FFmpeg](https://ffmpeg.org/) installed and accessible via the command line.

---

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/jsaleix/video-merger.git
   cd video-merger
   ```

2. Build the application:
   ```bash
   go build -o video-merger
   ```

3. Ensure `ffmpeg` is installed and in your system's PATH:
   ```bash
   ffmpeg -version
   ```

---

## Usage

1. Place the `video-merger` executable in a directory with video files (`.mp4`, `.mov`, `.mkv`).
2. Run the application:
   ```bash
   ./video-merger
   ```

3. The following files will be generated:
   - `result.mkv`: The merged video.
   - `timecode.txt`: A file listing the start time and name of each video.

---

## Example

### Input Directory:
```
video1.mp4
video2.mov
video3.mkv
```

### Generated Files:
1. **Merged Video**: `result.mkv`
2. **Timecode File**: `timecode.txt` containing:
   ```
   Timecodes

   00:00:00 - video1.mp4
   00:03:45 - video2.mov
   00:07:30 - video3.mkv
   ```

---

## Notes

- Only files with `.mp4`, `.mov`, or `.mkv` extensions are included.
- The merged video file will use the Matroska format (`.mkv`) for compatibility.
- Ensure all video files have compatible codecs, or they will be re-encoded.

---
