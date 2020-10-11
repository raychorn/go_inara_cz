package main

// ================================
import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func exampleScrape(someURL string) {
  // Request the HTML page.
  res, err := http.Get(someURL)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

  // Find the review items
  doc.Find("[name=\"searchcommodity\"]").Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
    s.Find("option").Each(func(ii int, item *goquery.Selection) {
		value, _ := item.Attr("value")
    	fmt.Printf("** %s --> %s\n", item.Text(), value)
	})
    //fmt.Printf("** %s --> %s\n", info, data)
  })
}


// ScrapeCommodities ...
func ScrapeCommodities(url string) string {
	exampleScrape(url)
	return fmt.Sprintf("Hello, ScrapeCommodities --> (%s).", url)
}

// Scraper ...
func Scraper(url string) string {
	return fmt.Sprintf("Hello, Scraper --> (%s).", url)
}
