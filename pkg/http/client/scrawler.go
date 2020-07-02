package client

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//ScrapingHead obtains data from HTML <head> tag
func ScrapingHead(hostname string) (string, string) {

	doc, err := goquery.NewDocument(hostname)
	if err != nil {
		log.Fatal(err)
	}
	var title string
	var link string

	title = doc.Find("title").Contents().Text()

	doc.Find("link").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, icon link
		href, _ := s.Attr("href")
		rel, _ := s.Attr("rel")
		//fmt.Printf("Review %d: %s - %s\n", i, rel, href)
		if strings.Contains(rel, "icon") {
			link = href

			if !strings.Contains(link, "http") {
				link = hostname + link
			}
			return false
		}
		return true
	},
	)
	fmt.Println("scraped: ", title, link)
	return title, link
}
