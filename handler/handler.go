package handler

import (
	"fmt"
	"io"
	"log"
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

	// write initial url content to file
	if depth == 0 {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		if err := h.downloaderClient.WriteContent(h.destinationDirectory, data); err != nil {
			fmt.Printf("%v", err)
			return
		}
	}

	links := h.crawler.ReadLink(resp.Body)
	log.Println("Links = ", len(links))
	validLinks := h.crawler.FilterLinks(links)
	log.Println("Valid links ", len(validLinks))
	for _, validLink := range validLinks {
		if depth+1 < crawler.MaxDepth {
			h.startCrawling(validLink.Url, depth+1)
		}
	}

}

func (h *Handler) WebCrawler() {
	h.startCrawling(h.initialURL, 0)
}