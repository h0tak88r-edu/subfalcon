package runner

import (
	"log"
	"strings"
	"time"

	"github.com/cyinnove/subfalcon/config"
	"github.com/cyinnove/subfalcon/pkg/db"
	"github.com/cyinnove/subfalcon/pkg/discord"
	"github.com/cyinnove/subfalcon/pkg/sub88r"
	"github.com/cyinnove/subfalcon/pkg/takeover"
)

func Run() {
	cfg := config.GetConfig()

	// Initialize database
	db.InitDB(config.DbFile)

	// Process domains
	var domains []string
	if cfg.SingleDomain != "" {
		domains = []string{cfg.SingleDomain}
	} else {
		// Read domains from file logic here
		// domains = readDomainsFromFile(cfg.DomainList)
	}

	for {
		for _, domain := range domains {
			// Initialize subdomain scraper
			subber := &sub88r.Subber{
				Domain:  domain,
				Results: &sub88r.Results{},
			}

			// Scrape subdomains from all sources
			subber.RapidDNS()
			subber.HackerTarget()
			subber.Anubis()
			subber.UrlScan()
			subber.Otx()
			subber.CrtSh()

			// Get unique subdomains
			newSubdomains := getUniqueSubdomains(subber.GetAllSubdomains())

			// Get existing subdomains from database
			existingSubdomains := db.Getsubdomains(config.DbFile)

			// Find truly new subdomains
			actualNewSubdomains := findNewSubdomains(newSubdomains, existingSubdomains)

			// Add new subdomains to database
			if len(actualNewSubdomains) > 0 {
				db.AddSubdmomains(actualNewSubdomains, config.DbFile)

				// Send new subdomains to Discord if webhook is configured
				if cfg.Webhook != "" {
					if err := discord.SendNewSubdomains(cfg.Webhook, actualNewSubdomains); err != nil {
						log.Printf("Error sending new subdomains to Discord: %v", err)
					}
				}
			}

			// Check for subdomain takeover if enabled
			if cfg.CheckTakeover {
				var takeoverResults []*takeover.SubdomainInfo

				for _, subdomain := range newSubdomains {
					info, err := takeover.CheckTakeover(subdomain)
					if err != nil {
						log.Printf("Error checking takeover for %s: %v", subdomain, err)
						continue
					}

					// Only append if there's actually a vulnerability
					if info.IsVulnerable {
						takeoverResults = append(takeoverResults, info)
					}
				}

				// Only send to Discord if there are actual vulnerable subdomains
				if cfg.Webhook != "" && len(takeoverResults) > 0 {
					if err := discord.SendTakeoverResults(cfg.Webhook, takeoverResults); err != nil {
						log.Printf("Error sending takeover results to Discord: %v", err)
					}
				}
			}
		}

		// If monitoring is not enabled, break the loop
		if !cfg.Monitor {
			break
		}

		// Wait for the monitoring interval before next scan
		time.Sleep(config.MonitorInterval)
	}
}

func getUniqueSubdomains(subdomains []string) []string {
	uniqueMap := make(map[string]bool)
	var result []string

	for _, subdomain := range subdomains {
		subdomain = strings.TrimSpace(subdomain)
		if subdomain != "" && !uniqueMap[subdomain] {
			uniqueMap[subdomain] = true
			result = append(result, subdomain)
		}
	}

	return result
}

func findNewSubdomains(current, existing []string) []string {
	existingMap := make(map[string]bool)
	for _, sub := range existing {
		existingMap[sub] = true
	}

	var newSubs []string
	for _, sub := range current {
		if !existingMap[sub] {
			newSubs = append(newSubs, sub)
		}
	}

	return newSubs
}
