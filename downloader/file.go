package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(url string, output string) {
	if output == "" {
		output = filepath.Base(url)
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	defer resp.Body.Close()

	file, err := os.Create(output)
	if err != nil {
		fmt.Println("Error creating file:", err)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error saving file:", err)
	}

	fmt.Println("Download completed:", output)
}
