package downloader

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const ChunkSize = 2 * 1024 * 1024

type Chunk struct {
	Start int
	End   int
	Data  []byte
	Err   error
}

func DownloadFile(url string, output string) {
	if output == "" {
		output = "download_file"
	}

	// Get file size
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

	// Create empty file
	file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Channel to receive downloaded chunks
	chunkChan := make(chan Chunk, fileSize/ChunkSize+1)
	var wg sync.WaitGroup

	// Download chunks in parallel
	for start := 0; start < fileSize; start += ChunkSize {
		end := start + ChunkSize - 1
		if end >= fileSize {
			end = fileSize - 1
		}
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			chunk := downloadChunk(url, start, end)
			chunkChan <- chunk
		}(start, end)
	}

	// Close channel when all downloads are done
	go func() {
		wg.Wait()
		close(chunkChan)
	}()

	// Write chunks to file
	for chunk := range chunkChan {
		if chunk.Err != nil {
			fmt.Println("Error downloading chunk:", chunk.Err)
			continue
		}
		// Seek to correct position and write
		_, err := file.Seek(int64(chunk.Start), 0)
		if err != nil {
			fmt.Println("Error seeking in file:", err)
			continue
		}
		_, err = file.Write(chunk.Data)
		if err != nil {
			fmt.Println("Error writing chunk:", err)
		}
	}
	fmt.Println("Download completed:", output)
}

func downloadChunk(url string, start int, end int) Chunk {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Chunk{Err: err}
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Chunk{Err: err}
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Chunk{Err: err}
	}
	return Chunk{Start: start, End: end, Data: data}
}
