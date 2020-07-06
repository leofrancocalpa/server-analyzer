package client

import (
	"fmt"
	"strings"

	"github.com/badoux/goscraper"
)

// GetTitleAndIcon return page Title, Icon
func GetTitleAndIcon(hostname string) (string, string) {
	s, err := goscraper.Scrape(hostname, 5)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	icon := s.Preview.Icon
	if !strings.Contains(icon, "http") {
		icon = hostname + icon
	}
	return s.Preview.Title, icon
}
