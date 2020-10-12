package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func closeDataChannel (dataChan chan<- map[string]int) {
	data := make(commodityItems)
	data["--DONE--"] = -1
	dataChan <- data
}

// GoScrape ...
func GoScrape(someURL string, dataChan chan<- map[string]int) {
  // Request the HTML page.
  res, err := http.Get(someURL)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()
  if res.StatusCode != 200 {
    log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
  }

  defer closeDataChannel(dataChan)

  // Load the HTML document
  doc, err := goquery.NewDocumentFromReader(res.Body)
  if err != nil {
    log.Fatal(err)
  }

  // Find the review items
  countItems := 0
  doc.Find("[name=\"searchcommodity\"]").Each(func(i int, s *goquery.Selection) {
    // For each item found, get the band and title
    s.Find("option").Each(func(ii int, item *goquery.Selection) {
		value, _ := item.Attr("value")
		data := make(map[string]int)
		i, _ := strconv.Atoi(value)
		data[item.Text()] = i
		dataChan <- data
		countItems++
    	fmt.Printf("** (%d) %s --> %s\n", countItems, item.Text(), value)
	})
  })
}


type commodityItems = map[string]int

// Commodities ...
var Commodities commodityItems

//CommoditiesURL ...
var CommoditiesURL string = "https://inara.cz/galaxy-commodities/"

// NewCommodities ...
func NewCommodities() {
    Commodities= make(commodityItems)
    dataChan := make(chan map[string]int, 5)
    go GoScrape(CommoditiesURL, dataChan)

    isDone := false
    for i := 0; i < 100; i++ {
        go func() {
            select {
                case data, ok := <-dataChan:
                    if (!ok) {
                        isDone = true
                        return
                    }
                    for k,v := range data {
                        if (k == "--DONE--") || (v == -1) { // data["--DONE--"] = -1
                            isDone = true
                            return
                        }
                        fmt.Printf("%s --> %d\n", k, v)
                        Commodities[k] = v
                        fmt.Printf("Num items %d\n", len(Commodities))
                    }
            }
        }()
        if (isDone) {
            fmt.Printf("isDone = %t\n", isDone)
            break
        }
    }
    if (!isDone) {
        fmt.Printf("Loop ended. (%t)\n", isDone)
    }
    fmt.Printf("Num items %d\n", len(Commodities))
    //Commodities["Tritium"] = 10269
    //Commodities["AgronomicTreatment"] = 10268
}

// AddCommodityItem ...
func AddCommodityItem(k string, v int) {
    Commodities[k] = v
}

// CommoditiesAsString ...
func CommoditiesAsString() string {
    b := new(bytes.Buffer)
    fmt.Fprint(b, "[")
    n := len(Commodities)
    i := 0
    for key, value := range Commodities {
        fmt.Fprintf(b, "{\"%s\":%d}", key, value)
        if (i < n-1) {
            fmt.Fprint(b, ", ")
        }
        i++
    }
    fmt.Fprint(b, "]")
    return b.String()
}

// CommodityNameByValue ...
func CommodityNameByValue(theValue int) string {
    items := map[int]string{}

    for key, value := range Commodities {
        items[value] = key
    }

    theKey, ok := items[theValue]
    if (!ok) {
        theKey = ""
    }
    return theKey
}