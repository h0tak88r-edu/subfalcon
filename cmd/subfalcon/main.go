package main

import (
	"flag"
	"log"

	"github.com/cyinnove/subfalcon/config"
	"github.com/cyinnove/subfalcon/runner"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	config.PrintLogo()

	// Define flag variables
	var domainList string
	var webhook string
	var monitor bool
	var singleDomain string // New flag for a single domain
	var checkTakeover bool  // New flag for subdomain takeover

	// Parse flags
	flag.StringVar(&domainList, "l", "", "Specify a file containing a list of domains")
	flag.StringVar(&webhook, "wh", "", "Specify the Discord webhook URL")
	flag.BoolVar(&monitor, "m", false, "Enable subdomain monitoring")
	flag.StringVar(&singleDomain, "d", "", "Specify a single domain for processing")           // New flag
	flag.BoolVar(&checkTakeover, "sdt", false, "Check for subdomain takeover vulnerabilities") // New flag
	flag.Parse()

	// Set the configuration values
	config.SetConfig(domainList, webhook, monitor, singleDomain, checkTakeover)

	// Validate flags
	if err := config.ValidateFlags(); err != nil {
		log.Fatal(err)
	}

	// Pass webhook URL to runner
	if webhook != "" {
		log.Printf("Using webhook URL: %s", webhook[:40]+"...") // Log truncated webhook URL for verification
	}
	runner.Run()
}
