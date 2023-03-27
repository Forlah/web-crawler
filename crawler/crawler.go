package crawler

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"web-crawler/model"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const MaxDepth = 3

//go:generate mockgen -source=crawler.go -destination=../mocks/crawler_mock.go -package=mocks
type Crawler interface {
	SetInitialUrl(url string) *CrawlerHelper
	ReadLink(body io.Reader) []model.Link
	FilterLinks([]model.Link) []model.Link
}

type CrawlerHelper struct {
	initialURL string
}

func New() Crawler {
	return &CrawlerHelper{}
}

func (c *CrawlerHelper) makeLink(tag html.Token, text string) model.Link {
	link := model.Link{
		Text: strings.TrimSpace(text),
	}

	for i := range tag.Attr {
		// get all href
		if tag.Attr[i].Key == "href" {
			link.Url = strings.TrimSpace(tag.Attr[i].Val)
		}
	}

	return link
}

func (c *CrawlerHelper) SetInitialUrl(url string) *CrawlerHelper {
	c.initialURL = url
	return c
}

func (c *CrawlerHelper) Valid(link model.Link) bool {
	if len(link.Text) == 0 {
		return false
	}

	if len(link.Url) == 0 || strings.Contains(strings.ToLower(link.Url), "javascript") {
		return false
	}

	if _, err := url.Parse(link.Url); err != nil {
		return false
	}

	// if !strings.HasPrefix(link.Url, "http://") || !strings.HasPrefix(link.Url, "https://") {
	// 	return false
	// }

	// initial_url, err := url.Parse(c.initialURL)
	// if err != nil {
	// 	return false
	// }

	// if !strings.HasPrefix(link_url.Host, initial_url.Host) {
	// 	return false
	// }

	return true
}

func (c *CrawlerHelper) ReadLink(body io.Reader) []model.Link {
	page := html.NewTokenizer(body)
	links := []model.Link{}

	var start *html.Token
	var text string
	for {
		_ = page.Next()
		token := page.Token()
		if token.Type == html.ErrorToken {
			break
		}
		if start != nil && token.Type == html.TextToken {
			text = fmt.Sprintf("%s%s", text, token.Data)
		}

		if token.DataAtom == atom.A {
			switch token.Type {
			case html.StartTagToken:
				if len(token.Attr) > 0 {
					start = &token
				}
			case html.EndTagToken:
				if start == nil {
					log.Println("Invalid link end tag without start tag. Skipping ...")
					continue
				}
				link := c.makeLink(*start, text)
				links = append(links, link)
				start = nil
				text = ""
			}

		}
	}
	return links
}

func (c *CrawlerHelper) FilterLinks(links []model.Link) []model.Link {
	filteredLinks := []model.Link{}
	for _, link := range links {
		if c.Valid(link) {
			log.Println("Valid link found: ", link)
			filteredLinks = append(filteredLinks, link)
		}
	}
	return filteredLinks
}
