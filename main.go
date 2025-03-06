package main

import (
	"flag"
	"fmt"
	"github.com/Zaman-R/go-idm/downloader"
)

func main() {
	url := flag.String("url", "", "URL of the file/video to download")
	output := flag.String("output", "", "Output filename")
	downloadType := flag.String("type", "file", "Download type (file/video/audio)")
	format := flag.String("format", "best", "Video format (e.g., best, mp4)")

	flag.Parse()

	if *url == "" {
		fmt.Println("Usage: go run main.go --url <URL> --type <file/video/audio> --output <filename>")
		return
	}

	switch *downloadType {
	case "file":
		downloader.DownloadFile(*url, *output)
	case "video":
		downloader.DownloadVideo(*url, *output, *format)
	case "formats":
		//downloader.GetAvailableFormats(*url)
	default:
		fmt.Println("Invalid type. Use 'file', 'video', or 'formats'.")
	}
}
