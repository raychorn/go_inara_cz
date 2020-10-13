package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	//validator "github.com/asaskevich/govalidator"

	"github.com/PuerkitoBio/goquery"
)

// CommodityDataJSON ...
type CommodityDataJSON struct {
	ID int    `json:"id"`
	Name  string `json:"name"`
}

func importFromJSON(fname string) map[string]CommodityDataStruct {
    fmt.Printf("DEBUG: importFromJSON(%s)\n", fname)
    rawJSON, _ := ioutil.ReadFile(fname)
    data := make(map[string]CommodityDataStruct)
    _ = json.Unmarshal([]byte(rawJSON), &data)
    return data
}


func exportAsJSON(items map[string]int) {
    n := len(items)
    datas := make(map[string]CommodityDataJSON, n)

	for k,v := range items {
        datas[k] = CommodityDataJSON{Name: k, ID: v}
        datas[fmt.Sprint(v)] = datas[k]
	}
    jsonString, err := json.Marshal(datas)
    if err != nil {
        log.Fatal(err)
    }
    now := time.Now()
    tsNow := strings.ReplaceAll(strings.ReplaceAll(now.Format(time.RFC3339), "-", "_"), ":", "")
    fmt.Printf("DEBUG: tsNow = %s\n", tsNow)
    //fmt.Printf("DEBUG: jsonString-->:%s\n", jsonString)
    fname := fmt.Sprintf("./inaracz_commodities-%s.json", tsNow)
    fmt.Printf("DEBUG: fname = %s\n", fname)
    err = ioutil.WriteFile(fname, jsonString, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

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
	data := make(map[string]int)
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

func fileinfoObj(info os.FileInfo) map[string]interface{} {
    return map[string]interface{}{
            "Name":    info.Name(),
            "Size":    info.Size(),
            "Mode":    info.Mode(),
            "ModTime": info.ModTime(),
            "IsDir":   info.IsDir(),
        }
}

// GetCommoditiesDataFileNames ...
func GetCommoditiesDataFileNames(isVerbose bool) map[string]map[string]interface{} {
    var foundFiles = make(map[string]map[string]interface{})
    var re = regexp.MustCompile(`(?m)^inaracz_commodities-(?P<ts>\d+_\d+_\d+T\d+_\d+)\.json$`)

    cutoff := 24 * time.Hour * 7
    //cutoff = 1 * time.Second
    root := "."
    now := time.Now()
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        match := re.FindAllString(path, -1)
        if (len(match) > 0) {
            fInfo := fileinfoObj(info)
            if diff := now.Sub(info.ModTime()); diff <= cutoff {
                foundFiles[path] = fInfo
            } else {
                fmt.Printf("DEBUG: Deleting \"%s\"", path)
                err := os.Remove(path) 
                if err != nil { 
                    log.Fatal(err) 
                }                 
            }
        }
        return nil
    })
    if err != nil {
        panic(err)
    }
    if (isVerbose) {
        fmt.Println("DEBUG: BEGIN: (2)")
        for k, v := range foundFiles {
            json, _ := json.Marshal(v)
            fmt.Printf("DEBUG: %s --> %s\n", k, json)
        }
        fmt.Println("DEBUG: END!!! (2)")
        fmt.Println()
    }

    return foundFiles
}

func keysFromFilesInterface(data map[string]map[string]interface{}) []string {
    keys := make([]string, 0, len(data))
    for k := range data {
        keys = append(keys, k)
    }
    return keys
}


// DoesCommoditiesDataExist ...
func DoesCommoditiesDataExist() (bool, map[string]map[string]interface{}) {
    files := GetCommoditiesDataFileNames(false)
    return len(files) > 0, files
}

func dumpCommoditiesData(data map[string]CommodityDataStruct) {
    fmt.Println("DEBUG:  BEGIN: dumpCommoditiesData")
    for k,v := range data {
        json, _ := json.Marshal(v)
        fmt.Printf("%s --> %s\n", k, json)
    }
    fmt.Println("DEBUG:  END!!! dumpCommoditiesData")
}

// Commodities ...
var Commodities map[string]int
// CommoditiesByValue ...
var CommoditiesByValue map[int]string

//CommoditiesURL ...
var CommoditiesURL string = "https://inara.cz/galaxy-commodities/"


// CommodityDataStruct ...
type CommodityDataStruct struct {
	id int
	name  string
}


// GetCommoditiesData ...
func GetCommoditiesData(data map[string]CommodityDataStruct) {
    Commodities = make(map[string]int)
    CommoditiesByValue = make(map[int]string)
    for k,v := range data {
        if (AreAllDigits(k)) {
            CommoditiesByValue[v.id] = v.name
        } else {
            Commodities[v.name] = v.id
        }
        json, _ := json.Marshal(v)
        fmt.Printf("%s --> %s\n", k, json)
    }
}


// NewCommodities ...
func NewCommodities(isVerbose bool) {
    Commodities = make(map[string]int)
    hasFiles, files := DoesCommoditiesDataExist()
    if (hasFiles) {
        keys := keysFromFilesInterface(files)
        data := importFromJSON(keys[0])
        dumpCommoditiesData(data)
        GetCommoditiesData(data)
        fmt.Println("DEBUG: Found the data.")
    } else {
        dataChan := make(chan map[string]int, 1)
        go GoScrape(CommoditiesURL, isVerbose, dataChan)

        isDone := false
        for i := 0; !isDone; i++ {
            select {
                case data, ok := <- dataChan:
                    if (!ok) {
                        isDone = true
                        log.Fatal("ERROR: Could not retrieve the whole data-set for commodities.")
                        break
                    }
                    for k,v := range data {
                        if (k == "--DONE--") && (v == -1) { 
                            isDone = true
                            break
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
        fmt.Printf("DEBUG: (%t)\n", isDone)
        if (isDone) {
            fmt.Println("JSON :: BEGIN")
            exportAsJSON(Commodities)
            fmt.Println("JSON :: END!!!")
        }
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
    theKey, ok := CommoditiesByValue[theValue]
    if (!ok) {
        theKey = ""
    }
    return theKey
}