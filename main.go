package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const MaxDepth = 2

var destinationDir string

type Link struct {
	url   string
	text  string
	depth int
}

type HttpError struct {
	errormsg string
}

func NewLink(tag html.Token, text string, depth int) Link {
	link := Link{
		text:  strings.TrimSpace(text),
		depth: depth,
	}

	for i := range tag.Attr {
		// get all href
		if tag.Attr[i].Key == "href" {
			link.url = strings.TrimSpace(tag.Attr[i].Val)
		}
	}

	return link
}

func (l Link) String() string {
	spacer := strings.Repeat("\t", l.depth)
	return fmt.Sprintf("%s%s (%d) - %s", spacer, l.text, l.depth, l.url)
}

func (l Link) Valid() bool {
	if l.depth >= MaxDepth {
		return false
	}
	if len(l.text) == 0 {
		return false
	}

	if len(l.url) == 0 || strings.Contains(strings.ToLower(l.url), "javascript") {
		return false
	}

	// Todo: check against initial url
	return true

}

func (httpErr HttpError) Error() string {
	return httpErr.errormsg
}

func LinkReader(resp *http.Response, depth int) []Link {
	page := html.NewTokenizer(resp.Body)
	links := []Link{}

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
					//log.Println("Invalid link end tag without start tag. Skipping ...")
					continue
				}
				link := NewLink(*start, text, depth)
				if link.Valid() {
					links = append(links, link)
					//log.Println("Valid link found: ", link)
				}
				start = nil
				text = ""
			}

		}
	}
	return links
}

func downloader(url string) (*http.Response, error) {
	//fmt.Printf("Downloading ... %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error: %v", err)
	}

	//defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return nil, HttpError{errormsg: fmt.Sprintf("Error (%d): %s", resp.StatusCode, url)}
	}

	return resp, nil

}

func recursiveDownloader(url string, depth int) {
	page, err := downloader(url)
	if err != nil {
		return
	}

	// put main page content in file
	if depth == 0 {
		f, err := os.Create(destinationDir)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		data, err := io.ReadAll(page.Body)
		if err != nil {
			log.Fatal(err)
		}
		n, err := f.WriteString(string(data))
		if err != nil {
			fmt.Println("Error copying content to file", err.Error())
			return
		}
		fmt.Println("Wrritten ", n)
	}

	links := LinkReader(page, depth)
	for _, link := range links {
		if depth+1 < MaxDepth {
			recursiveDownloader(link.url, depth+1)
		}
	}

}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Please enter a starting URL: ")
	startingURL, err := reader.ReadString('\n')
	if err != nil {
		panic("unable to read starting url")
	}

	fmt.Print("Please enter your destination directory: ")
	destinationDir, err = reader.ReadString('\n')
	if err != nil {
		panic("unable to read destination directory")
	}

	recursiveDownloader(strings.TrimSpace(startingURL), 0)

}
