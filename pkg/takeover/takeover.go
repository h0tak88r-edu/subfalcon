package takeover

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
)

type SubdomainInfo struct {
	Subdomain    string
	CNAME        string
	Status       string
	IsVulnerable bool
}

const (
	Red   = "\033[0;31m"
	Blue  = "\033[0;34m"
	Green = "\033[0;32m"
	NC    = "\033[0m" // No Color
)

func CheckTakeover(subdomain string) (*SubdomainInfo, error) {
	// Get CNAME record
	cname, err := net.LookupCNAME(subdomain)
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			fmt.Printf("%s[NXDOMAIN] %s%s\n", Green, subdomain, NC)
			return &SubdomainInfo{
				Subdomain: subdomain,
				Status:    "NXDOMAIN",
			}, nil
		}
		return nil, err
	}

	info := &SubdomainInfo{
		Subdomain: subdomain,
		CNAME:     strings.TrimSuffix(cname, "."),
		Status:    "Active",
	}

	// Check for Azure services takeover
	if isAzureVulnerable(info) {
		info.IsVulnerable = true
		info.Status = "Potentially Vulnerable (Azure)"
		fmt.Printf("%s[VULNERABLE] %s -> CNAME: %s%s\n", Red, subdomain, info.CNAME, NC)
	} else {
		fmt.Printf("%s[ACTIVE] %s -> CNAME: %s%s\n", Blue, subdomain, info.CNAME, NC)
	}

	return info, nil
}

func isAzureVulnerable(info *SubdomainInfo) bool {
	azureRegex := regexp.MustCompile(`(?i)^(?:[a-z0-9-]+\.)?(?:cloudapp\.net|azurewebsites\.net|cloudapp\.azure\.com)$`)

	if info.Status == "NXDOMAIN" && azureRegex.MatchString(info.CNAME) {
		url := fmt.Sprintf("https://%s", info.CNAME)
		_, err := http.Get(url)
		return err != nil
	}
	return false
}
