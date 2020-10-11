package main

import (
	"fmt"
	"strings"
	scraper "webscraper"

	"github.com/spf13/cobra"
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
    Use:   "scrape [url string]",
    Short: "Scrape using a url - must reference inara.cz",
    Long: `scrape an inara.cz url.`,
    Args: cobra.MinimumNArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
      url := strings.Join(args, " ")
      fmt.Println("Scrape: " + url)
      fmt.Println(scraper.Scraper(url))
    },
  }

  // go run ./go_inara_cz.go scrape "url"
  
  var rootCmd = &cobra.Command{Use: "go_inara_cz"}
  rootCmd.AddCommand(cmdPrint, cmdEcho)
  rootCmd.AddCommand(cmdScrape)
  rootCmd.Execute()
}
