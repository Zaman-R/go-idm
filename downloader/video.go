package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Chunk size (in bytes)
const chunkSize = 1024 * 1024 // 5 MB per chunk

// DownloadChunk downloads a specific chunk of a file
func DownloadChunk(url string, start int64, end int64, filePath string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Prepare HTTP request with Range header to fetch the chunk
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error downloading chunk:", err)
		return
	}
	defer resp.Body.Close()

	// Open file for writing
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Move the pointer to the right position for the chunk
	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		fmt.Println("Error seeking file:", err)
		return
	}

	// Write the chunk to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error saving chunk:", err)
		return
	}

	fmt.Printf("Downloaded chunk %d-%d\n", start, end)
}

// DownloadVideo downloads a video by splitting it into chunks
func DownloadVideo(URL string, output string, format string) {
	if output == "" {
		output = "video.mp4"
	}

	// Get video info (length of video file)
	cmd := exec.Command("yt-dlp", "-f", format, "-g", URL)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error getting video URL:", err)
		return
	}

	// Extract direct video URL and clean it
	videoURL := strings.TrimSpace(string(out))

	// If we encounter a warning, we want to extract only the URL part
	if strings.Contains(videoURL, "WARNING") {
		fmt.Println("Warning in yt-dlp output, cleaning up URL.")

		// Check if the output contains more than one line, otherwise, just use the output as is
		lines := strings.Split(videoURL, "\n")
		if len(lines) > 1 {
			videoURL = lines[len(lines)-1] // Get the last line (URL)
		} else {
			// If the output is just one line and contains a URL, use it directly
			videoURL = lines[0]
		}
	}

	// Decode the URL
	decodedURL, err := url.QueryUnescape(videoURL)
	if err != nil {
		fmt.Println("Error unescaping URL:", err)
		return
	}

	// Ensure the URL is well-formed
	_, err = url.ParseRequestURI(decodedURL)
	if err != nil {
		fmt.Println("Error parsing video URL:", err)
		return
	}

	// Get file size (to split into chunks)
	resp, err := http.Head(decodedURL)
	if err != nil {
		fmt.Println("Error getting file size:", err)
		return
	}
	defer resp.Body.Close()

	// Get content length (total size of the video)
	fileSize := resp.ContentLength
	fmt.Println("File size:", fileSize)

	// Create output file
	file, err := os.Create(output)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	file.Close()

	// Calculate the number of chunks (use int64 for consistency)
	numChunks := fileSize / chunkSize
	if fileSize%chunkSize != 0 {
		numChunks++
	}

	// Download the video in chunks concurrently
	var wg sync.WaitGroup
	for i := int64(0); i < numChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize - 1
		if end > fileSize-1 {
			end = fileSize - 1
		}

		wg.Add(1)
		go DownloadChunk(decodedURL, start, end, output, &wg)
	}

	// Wait for all chunks to download
	wg.Wait()

	fmt.Println("Video download completed:", output)
}
