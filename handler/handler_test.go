package handler

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"web-crawler/mocks"
	"web-crawler/model"

	"github.com/golang/mock/gomock"
)

func TestWebCrawlerHandler(t *testing.T) {
	t.Run("Test web crawler successfully.", func(t *testing.T) {
		controller := gomock.NewController(t)
		defer controller.Finish()
		startUrl := "https://google.com"
		destinationDir := "../crawler/sample.txt"

		mockCrawler := mocks.NewMockCrawler(controller)
		mockDownloaderClient := mocks.NewMockDownloaderClient(controller)

		handler := New(startUrl, destinationDir)
		handler.crawler = mockCrawler
		handler.downloaderClient = mockDownloaderClient

		httpRespMock := http.Response{
			Body: io.NopCloser(bytes.NewBufferString("Hello Trivity")),
		}

		mockDownloaderClient.EXPECT().DownloadURL(startUrl).Return(&httpRespMock, nil)
		mockDownloaderClient.EXPECT().WriteContent(destinationDir, gomock.Any()).Return(nil)

		mockCrawler.EXPECT().SetInitialUrl(startUrl)
		mockLinks := []model.Link{
			{
				Url: "https://es.com",
			},
		}

		mockCrawler.EXPECT().ReadLink(gomock.All()).Return(mockLinks).AnyTimes()
		mockCrawler.EXPECT().FilterLinks(gomock.All()).Return(mockLinks).AnyTimes()
		mockDownloaderClient.EXPECT().DownloadURL(gomock.Any()).Return(&httpRespMock, nil).AnyTimes()

		handler.WebCrawler()
	})
}
