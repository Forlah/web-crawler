package crawler

import (
	"bytes"
	"os"
	"testing"
	"web-crawler/model"

	"github.com/stretchr/testify/assert"
)

func Test_ReadLink(t *testing.T) {
	t.Run("Test Reading web page content successfully", func(t *testing.T) {
		content, err := os.ReadFile("sample.txt")
		assert.NoError(t, err)

		crawler := New()
		links := crawler.ReadLink(bytes.NewBuffer(content))
		assert.NotEmpty(t, links)
	})

}

func Test_FilterLinks(t *testing.T) {
	t.Run("Test filter links", func(t *testing.T) {
		crawler := New()
		links := []model.Link{
			{
				Url:  "https://eu.wikipedia.com",
				Text: "EU",
			},
		}
		result := crawler.FilterLinks(links)
		assert.Equal(t, int(1), len(result))
	})
}
