package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const ChunkSize = 1024 * 1024

func DownloadFile(url string, output string) {
	if output == "" {
		output = "download_file"
	}

	resp, err := http.Head(url)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fileSize, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		fmt.Println("Failed to get file size:", err)
		return
	}
	fmt.Printf("File size: %d bytes\n", fileSize)

	file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	for start := 0; start < fileSize; start += ChunkSize {
		end := start + ChunkSize - 1
		if end >= fileSize {
			end = fileSize - 1
		}
		fmt.Printf("Downloading chunk: %d-%d\n", start, end)
		if err := downloadChunk(url, file, start, end); err != nil {
			fmt.Println("Error downloading chunk:", err)
			return
		}
	}
	fmt.Println("Download completed:", output)
}

func downloadChunk(url string, file *os.File, start, end int) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = file.Seek(int64(start), 0)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	return err
}
