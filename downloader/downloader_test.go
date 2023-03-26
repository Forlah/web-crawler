package downloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DownloadURL(t *testing.T) {
	t.Run("Test download from url successfully", func(t *testing.T) {
		downloaderClient := New()
		resp, err := downloaderClient.DownloadURL("https://google.com")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
	})
}

func Test_WriteContent(t *testing.T) {
	t.Run("Test write content to file successfully", func(t *testing.T) {
		downloaderClient := New()
		err := downloaderClient.WriteContent("text.txt", []byte("Hello Trivity"))
		assert.NoError(t, err)
	})
}
