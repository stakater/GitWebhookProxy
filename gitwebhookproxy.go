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
	secret        = flagSet.String("secret", "", "Secret of the Webhook API. If not set validation is not made.")
	provider      = flagSet.String("provider", "github", "Git Provider which generates the Webhook")
	allowedPaths  = flagSet.String("allowedPaths", "", "Comma-Separated String List of allowed paths")
	ignoredUsers  = flagSet.String("ignoredUsers", "", "Comma-Separated String List of users to ignore while proxying Webhook request")
	allowedUsers  = flagSet.String("allowedUser", "", "Comma-Separated String List of users to allow while proxying Webhook request")
)

func validateRequiredFlags() {
	isValid := true
	if len(strings.TrimSpace(*upstreamURL)) == 0 {
		log.Println("Required flag 'upstreamURL' not specified")
		isValid = false
	}

	if !isValid {
		fmt.Println("")
		//TODO: Usage not working as expected in flagSet
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
	allowedPathsArray := []string{}
	if len(*allowedPaths) > 0 {
		allowedPathsArray = strings.Split(*allowedPaths, ",")
	}

	// Split Comma-Separated list into an array
	ignoredUsersArray := []string{}
	if len(*ignoredUsers) > 0 {
		ignoredUsersArray = strings.Split(*ignoredUsers, ",")
	}

	log.Printf("Stakater Git WebHook Proxy started with provider '%s'\n", lowerProvider)
	p, err := proxy.NewProxy(*upstreamURL, allowedPathsArray, lowerProvider, *secret, ignoredUsersArray)
	if err != nil {
		log.Fatal(err)
	}

	if err := p.Run(*listenAddress); err != nil {
		log.Fatal(err)
	}

}
