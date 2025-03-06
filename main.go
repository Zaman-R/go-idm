package main

import (
	"fmt"
	"github.com/Zaman-R/go-idm/cli"
	"log"
)

func main() {
	options := cli.ParseArgs()
	if options == nil {
		return
	}

	fmt.Println("Downloading %s as %s...\n", options.URL, options.MediaType)

	switch options.MediaType {
	case "file":
		downloader.DownloadFile(options.URL, options.Output)
	case "video":
		downloader.DownloadVideo(options.URL, options.Output)
	case "audio":
		downloader.DownloadAudio(options.URL, options.Output)
	default:
		log.Println("Invalid type. Use: file, video, or audio")
	}
}
