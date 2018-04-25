package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/namsral/flag"
	"github.com/stakater/GitWebhookProxy/pkg/proxy"
)

var (
	flagSet       = flag.NewFlagSetWithEnvPrefix(os.Args[0], "GWP", 0)
	listenAddress = flagSet.String("listen", ":8080", "Address on which the proxy listens.")
	upstreamURL   = flagSet.String("upstreamURL", "", "URL to which the proxy requests will be forwarded (required)")
	secret        = flagSet.String("secret", "", "Secret of the Webhook API (required)")
	provider      = flagSet.String("provider", "github", "Git Provider which generates the Webhook")
	allowedPaths  = flagSet.String("allowedPaths", "", "Comma-Separated String List of allowed paths")
)

func validateRequiredFlags() {
	isValid := true
	if len(strings.TrimSpace(*upstreamURL)) == 0 {
		log.Println("Required flag 'upstreamURL' not specified")
		isValid = false
	}
	if len(strings.TrimSpace(*secret)) == 0 {
		log.Println("Required flag 'secret' not specified")
		isValid = false
	}

	if !isValid {
		fmt.Println("")
		flagSet.Usage()
		fmt.Println("")

		panic("See Flag Usage")
	}
}

func main() {
	flagSet.Parse(os.Args[1:])
	validateRequiredFlags()
	lowerProvider := strings.ToLower(*provider)

	// Split Comma-Separated list into an array
	var allowedPathsArray []string
	if len(*allowedPaths) > 0 {
		allowedPathsArray = strings.Split(*allowedPaths, ",")
	}

	log.Printf("Stakater Git WebHook Proxy started with provider '%s'\n", lowerProvider)
	proxy.NewProxy(*listenAddress, *upstreamURL, allowedPathsArray, lowerProvider, *secret)
}
