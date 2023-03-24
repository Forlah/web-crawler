package downloader

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type DownloaderClient interface {
	DownloadURL(url string) (*http.Response, error)
	WriteContent(filepath string, data []byte) error
}

type downloader struct{}

func New() DownloaderClient {
	return &downloader{}
}

func (d *downloader) DownloadURL(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Error: %v", err))
	}

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf(fmt.Sprintf("Error (%d): %s", resp.StatusCode, url))
	}

	return resp, nil
}

func (d *downloader) WriteContent(filepath string, data []byte) error {
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	n, err := file.WriteString(string(data))
	if err != nil {
		fmt.Println("Error writing content to file", err.Error())
		return err
	}

	fmt.Printf("Wrote %d bytes", n)
	return nil
}
