package scraper

import (
	"fmt"
)

// ScrapeCommodities ...
func ScrapeCommodities(url string) string {
	return fmt.Sprintf("Hello, ScrapeCommodities --> (%s).", url)
}

// Scraper ...
func Scraper(url string) string {
	return fmt.Sprintf("Hello, Scraper --> (%s).", url)
}
