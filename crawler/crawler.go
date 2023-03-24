package crawler

import (
	"fmt"
	"io"
	"log"
	"strings"
	"web-crawler/model"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const MaxDepth = 3

type Crawler interface {
	ReadLink(body io.Reader) []model.Link
	FilterLinks([]model.Link) []model.Link
}

type crawler struct{}

func New() Crawler {
	return &crawler{}
}

func (c crawler) makeLink(tag html.Token, text string) model.Link {
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

func (c crawler) Valid(link model.Link) bool {
	if link.Depth >= MaxDepth {
		return false
	}
	if len(link.Text) == 0 {
		return false
	}

	if len(link.Url) == 0 || strings.Contains(strings.ToLower(link.Url), "javascript") {
		return false
	}

	// Todo: check against initial url
	return true

}

func (c crawler) ReadLink(body io.Reader) []model.Link {
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

func (c crawler) FilterLinks(links []model.Link) []model.Link {
	filteredLinks := []model.Link{}
	for _, link := range links {
		if c.Valid(link) {
			log.Println("Valid link found: ", link)
			filteredLinks = append(filteredLinks, link)
		}
	}
	return filteredLinks
}
