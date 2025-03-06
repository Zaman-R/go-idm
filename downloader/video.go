package downloader

import (
	"fmt"
	"os/exec"
	"strings"
)

func DownloadVideo(url string, output string, format string) {
	if output == "" {
		output = "video.mp4"
	}

	// Build yt-dlp command
	cmd := exec.Command("yt-dlp", "-f", format, "-o", output, url)

	// Capture Output
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println(string(out))
	fmt.Println("Video Downloaded Successfully", output)

}

// GetAvailableFormats lists available formats for a video
func GetAvailableFormats(url string) {
	cmd := exec.Command("yt-dlp", "-F", url)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Println("Available formats:\n", string(out))
}

// ExtractDirectURL gets the direct video URL (for chunk-based downloads)
func ExtractDirectURL(url string, format string) (string, error) {
	cmd := exec.Command("yt-dlp", "-f", format, "-g", url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// Trim spaces & return direct URL
	return strings.TrimSpace(string(out)), nil
}
