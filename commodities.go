package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// TimeoutDialer ...
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
    return func(netw, addr string) (net.Conn, error) {
        conn, err := net.DialTimeout(netw, addr, cTimeout)
        if err != nil {
            return nil, err
        }
        conn.SetDeadline(time.Now().Add(rwTimeout))
        return conn, nil
    }
}

// NewTimeoutClient ...
func NewTimeoutClient(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Client {

    return &http.Client{
        Transport: &http.Transport{
            Dial: TimeoutDialer(connectTimeout, readWriteTimeout),
        },
    }
}

func closeDataChannel (dataChan chan<- map[string]int) {
	data := make(commodityItems)
	data["--DONE--"] = -1
	dataChan <- data
}

// GoScrape ...
func GoScrape(someURL string, isVerbose bool, dataChan chan<- map[string]int) {
    // Request the HTML page.
    client := NewTimeoutClient(time.Duration(1000*1000*1000*30), time.Duration(1000*1000*1000*30))
    res, err := client.Get(someURL)
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

    if (isVerbose) {
        fmt.Println("GoScrape :: Signal Start.")
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
            if (isVerbose) {
                countItems++
                fmt.Printf("GoScrape :: (%d) %s --> %s\n", countItems, item.Text(), value)
            }
        })
    })
}


type commodityItems = map[string]int

// Commodities ...
var Commodities commodityItems

//CommoditiesURL ...
var CommoditiesURL string = "https://inara.cz/galaxy-commodities/"

// NewCommodities ...
func NewCommodities(isVerbose bool) {
    Commodities= make(commodityItems)
    dataChan := make(chan map[string]int, 1)
    go GoScrape(CommoditiesURL, isVerbose, dataChan)

    isDone := false
    for i := 0; !isDone; i++ {
        select {
            case data, ok := <- dataChan:
                if (!ok) {
                    isDone = true
                    return
                }
                for k,v := range data {
                    if (k == "--DONE--") && (v == -1) { 
                        isDone = true
                        return
                    }
                    if (isVerbose) {
                        fmt.Printf("NewCommodities :: %s --> %d\n", k, v)
                    }
                    Commodities[k] = v
                    if (isVerbose) {
                        fmt.Printf("NewCommodities :: Num items %d\n", len(Commodities))
                    }
                }
        }
        if (isDone) {
            if (isVerbose) {
                fmt.Printf("NewCommodities :: isDone = %t\n", isDone)
            }
            break
        }
    }
    if (!isDone) {
        if (isVerbose) {
            fmt.Printf("NewCommodities :: Loop ended. (%t)\n", isDone)
        }
    }
    if (isVerbose) {
        fmt.Printf("NewCommodities :: Num items %d\n", len(Commodities))
    }
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