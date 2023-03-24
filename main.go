package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"web-crawler/handler"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Please enter a starting URL: ")
	startingURL, err := reader.ReadString('\n')
	if err != nil {
		panic("unable to read starting url")
	}

	fmt.Print("Please enter your destination directory: ")
	destinationDir, err := reader.ReadString('\n')
	if err != nil {
		panic("unable to read destination directory")
	}

	handler := handler.New(strings.TrimSpace(startingURL), strings.TrimSpace(destinationDir))
	handler.WebCrawler()

}
