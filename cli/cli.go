package cli

import (
	"flag"
	"fmt"
)

type Options struct {
	URL       string
	Output    string
	MediaType string
}

func ParseArgs() *Options {
	url := flag.String("url", "", "URL of the file/video/audio to download")
	output := flag.String("output", "", "Output file name (optional)")
	mediaType := flag.String("type", "file", "Specify type: file, video, or audio")

	flag.Parse()

	if *url == "" {
		fmt.Println("URL flag is required")
		flag.Usage()
		return nil
	}

	return &Options{
		URL:       *url,
		Output:    *output,
		MediaType: *mediaType,
	}
}
