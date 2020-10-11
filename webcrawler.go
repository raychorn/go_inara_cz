package main

// ================================
import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over token attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	
	return
}

func getInfo(t html.Token) (ok bool, href string) {
	// Iterate over token attributes until we find an "href"
	fmt.Println("BEGIN: Attrs")
	for _, a := range t.Attr {
		fmt.Printf("%s --> %s", a.Key, a.Val)
	}
	fmt.Println("END!!! Attrs")
	
	return
}

// Extract all http** links from a given webpage
func crawl(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()

	if err != nil {
		fmt.Println("ERROR: Failed to crawl:", url)
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function completes

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			isSelect := t.Data == "select"
			if !isAnchor && !isSelect {
				continue
			}

			// Extract the href value, if there is one
			if (isAnchor) {
				ok, url := getHref(t)
				if !ok {
					continue
				}

				// Make sure the url begines in http**
				hasProto := strings.Index(url, "http") == 0
				if hasProto {
					ch <- url
				}
				continue
			}

			if (isSelect) {
				fmt.Println("Found Select.")
				t := z.Token()
				ok, url := getInfo(t)
				fmt.Printf("ok is %t, url is %s", ok, url)
			}
			continue
		default:
			//fmt.Println(tt)
			continue
		}
	}
}

func doIt(theURL string) {
	foundUrls := make(map[string]bool)
	seedUrls := []string{theURL}

	maxDepth := 1
	currentDepth := 1

	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool) 

	// Kick off the crawl process (concurrently)
	for _, url := range seedUrls {
		go crawl(url, chUrls, chFinished)
		break
	}

	// Subscribe to both channels
	for c := 0; c < len(seedUrls); {
		select {
			case url := <-chUrls:
				foundUrls[url] = true
				if (currentDepth >= maxDepth) {
					c++
				}
				currentDepth++
			case <-chFinished:
				c++
		}
	}

	// We're done! Print the results...

	fmt.Println("\nFound", len(foundUrls), "unique urls:")

	for url := range foundUrls {
		fmt.Println(" - " + url)
	}

	close(chUrls)
}
