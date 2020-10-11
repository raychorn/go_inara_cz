package scraper

import (
	"fmt"
)

// Scraper ...
func Scraper(url string) string {
	return fmt.Sprintf("Hello world --> (%s).", url)
}
