package main

import (
	"fmt"
	"strings"

	cobra "github.com/spf13/cobra"
)

func main() {
  var cmdPrint = &cobra.Command{
    Use:   "print [string to print]",
    Short: "Print anything to the screen",
    Long: `print is for printing anything back to the screen.
For many years people have printed back to the screen.`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("Print: " + strings.Join(args, " "))
    },
  }

  var cmdEcho = &cobra.Command{
    Use:   "echo [string to echo]",
    Short: "Echo anything to the screen",
    Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("Echo: " + strings.Join(args, " "))
    },
  }

  var cmdScrape = &cobra.Command{
    Use:   "scrape",
    Short: "Scrape via inara.cz",
    Long: `scrape an inara.cz url using the standard url.`,
    Args: cobra.MinimumNArgs(0),
    Run: func(cmd *cobra.Command, args []string) {
      //url := strings.Join(args, " ")
      NewCommodities(false)
      cItems := Commodities
      //fmt.Printf("cItems = %s\n", CommoditiesAsString())
      tritium := cItems["Tritium"]
      cName := CommodityNameByValue(tritium)
      fmt.Printf("%s --> %d\n", cName, tritium)
      fmt.Printf("Scrape: %d items\n", len(cItems))
    },
  }

  var cmdScrapeCommodities = &cobra.Command{
    Use:   "scrapecommodities [url string]",
    Short: "Scrape Commodities using default built-in",
    Long: `scrape commodities via inara.cz url.`,
    Args: cobra.MinimumNArgs(0),
    Run: func(cmd *cobra.Command, args []string) {
      if (len(args) > 0) {
        fmt.Println("No need to use an argument for this command.")
      }
      url := CommoditiesURL
      fmt.Println("Scrape Commodities: " + url)
      fmt.Println(ScrapeCommodities(url))
    },
  }

  // go run ./go_inara_cz.go scrape "url"
  // commodities_buymin_url = "https://inara.cz/ajaxaction.php?act=goodsdata&refname=buymin&refid={}&refid2={}".format(commodity_refid, star_system_refid)
  // commodities_sellmax_url = "https://inara.cz/ajaxaction.php?act=goodsdata&refname=sellmax&refid={}&refid2={}".format(commodity_refid, star_system_refid)


  var rootCmd = &cobra.Command{Use: "go_inara_cz"}
  rootCmd.AddCommand(cmdPrint, cmdEcho)
  rootCmd.AddCommand(cmdScrape)
  rootCmd.AddCommand(cmdScrapeCommodities)
  rootCmd.Execute()
}
