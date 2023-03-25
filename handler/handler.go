package handler

import (
	"fmt"
	"io"
	"log"
	"sync"
	"web-crawler/crawler"
	"web-crawler/downloader"
)

type Handler struct {
	crawler              crawler.Crawler
	downloaderClient     downloader.DownloaderClient
	initialURL           string
	destinationDirectory string
}

func New(startingURL, destinationDir string) *Handler {
	return &Handler{
		crawler:              crawler.New(),
		downloaderClient:     downloader.New(),
		initialURL:           startingURL,
		destinationDirectory: destinationDir,
	}
}

func (h *Handler) startCrawling(url string, depth int) {
	resp, err := h.downloaderClient.DownloadURL(url)
	if err != nil {
		fmt.Printf("downloading url, Error: %v", err)
		return
	}

	links := h.crawler.ReadLink(resp.Body)
	validLinks := h.crawler.FilterLinks(links)
	fmt.Println("Valid links count ", len(validLinks))
	for _, validLink := range validLinks {
		if depth+1 < crawler.MaxDepth {
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				h.startCrawling(validLink.Url, depth+1)
			}()

			// Wait for function to be executed
			wg.Wait()
		}
	}

}

func (h *Handler) WebCrawler() {
	resp, err := h.downloaderClient.DownloadURL(h.initialURL)
	if err != nil {
		fmt.Printf("downloading initial url, Error: %v", err)
		return
	}

	// write initial url content to file
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if err := h.downloaderClient.WriteContent(h.destinationDirectory, data); err != nil {
		fmt.Printf("%v", err)
		return
	}

	h.crawler.SetInitialUrl(h.initialURL)
	// recursively crawl webpage
	h.startCrawling(h.initialURL, 0)
}
