package main

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/stakater/GitWebhookProxy/pkg/proxy"
)

// FIPS mode is enabled by setting the GOFIPS140 environment variable
// To build with FIPS mode: export GOFIPS140='v1.0.0' && go build -o gitwebhookproxy-fips examples/main.go

func main() {
	// Define command-line flags
	listenAddress := pflag.String("listen", ":8080", "Address on which the proxy listens.")
	upstreamURL := pflag.String("upstreamURL", "http://localhost:8081", "URL to which the proxy requests will be forwarded")
	secret := pflag.String("secret", "", "Secret of the Webhook API. If not set validation is not made.")
	provider := pflag.String("provider", "github", "Git Provider which generates the Webhook (github or gitlab)")
	allowedPaths := pflag.String("allowedPaths", "", "Comma-Separated String List of allowed paths")
	ignoredUsers := pflag.String("ignoredUsers", "", "Comma-Separated String List of users to ignore while proxying Webhook request")
	allowedUsers := pflag.String("allowedUsers", "", "Comma-Separated String List of users to allow while proxying Webhook request")

	// Parse flags
	pflag.Parse()

	// Validate required flags
	if *upstreamURL == "" {
		log.Println("Required flag 'upstreamURL' not specified")
		pflag.Usage()
		os.Exit(1)
	}

	// Parse comma-separated lists
	allowedPathsArray := []string{}
	if *allowedPaths != "" {
		allowedPathsArray = splitCommaSeparatedList(*allowedPaths)
	}

	ignoredUsersArray := []string{}
	if *ignoredUsers != "" {
		ignoredUsersArray = splitCommaSeparatedList(*ignoredUsers)
	}

	allowedUsersArray := []string{}
	if *allowedUsers != "" {
		allowedUsersArray = splitCommaSeparatedList(*allowedUsers)
	}

	// Create and run the proxy
	log.Printf("Starting Git WebHook Proxy with provider '%s'", *provider)

	// When GOFIPS140 environment variable is set, FIPS mode will be automatically enabled
	// via the init() function in fips_config.go
	p, err := proxy.NewProxy(*upstreamURL, allowedPathsArray, *provider, *secret, ignoredUsersArray, allowedUsersArray)
	if err != nil {
		log.Fatal(err)
	}

	if err := p.Run(*listenAddress); err != nil {
		log.Fatal(err)
	}
}

// Helper function to split comma-separated list
func splitCommaSeparatedList(list string) []string {
	if list == "" {
		return []string{}
	}

	// Simple string split by comma
	items := strings.Split(list, ",")

	// Trim spaces from each item
	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}

	return items
}
