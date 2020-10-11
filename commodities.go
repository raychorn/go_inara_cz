package main

import (
	"bytes"
	"fmt"
)

type commodityItems = map[string]int

// Commodities ...
var Commodities commodityItems

// NewCommodities ...
func NewCommodities(){
	Commodities= make(commodityItems)
    Commodities["Tritium"] = 10269
    Commodities["AgronomicTreatment"] = 10268
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

