package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/stakater/GitWebhookProxy/pkg/proxy"
)

var (
	listenAddress = pflag.String("listen", ":8080", "Address on which the proxy listens.")
	upstreamURL   = pflag.String("upstreamURL", "", "URL to which the proxy requests will be forwarded (required)")
	secret        = pflag.String("secret", "", "Secret of the Webhook API. If not set validation is not made.")
	provider      = pflag.String("provider", "github", "Git Provider which generates the Webhook")
	allowedPaths  = pflag.String("allowedPaths", "", "Comma-Separated String List of allowed paths")
	ignoredUsers  = pflag.String("ignoredUsers", "", "Comma-Separated String List of users to ignore while proxying Webhook request")
	allowedUsers  = pflag.String("allowedUser", "", "Comma-Separated String List of users to allow while proxying Webhook request")
	insecureTLS   = pflag.Bool("insecureSkipVerify", false, "Skip TLS certificate verification of the upstream. INSECURE: only enable for trusted upstreams using self-signed certificates.")
)

func validateRequiredFlags() {
	isValid := true
	if len(strings.TrimSpace(*upstreamURL)) == 0 {
		log.Println("Required flag 'upstreamURL' not specified")
		isValid = false
	}

	if !isValid {
		fmt.Println("")
		pflag.Usage()
		fmt.Println("")

		panic("See Flag Usage")
	}
}

func main() {
	// First parse command line flags
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.Parse()

	// Then process environment variables, but only if flags weren't set via command line
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "GWP_") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) != 2 {
				continue
			}

			// Map environment variables to flag names, preserving case
			envName := strings.TrimPrefix(parts[0], "GWP_")
			name := ""
			switch strings.ToUpper(envName) {
			case "UPSTREAMURL":
				name = "upstreamURL"
			case "ALLOWEDPATHS":
				name = "allowedPaths"
			case "IGNOREDUSERS":
				name = "ignoredUsers"
			case "ALLOWEDUSER":
				name = "allowedUser"
			case "PROVIDER":
				name = "provider"
			case "SECRET":
				name = "secret"
			case "LISTEN":
				name = "listen"
			case "INSECURESKIPVERIFY":
				name = "insecureSkipVerify"
			default:
				name = strings.ToLower(envName)
			}
			value := parts[1]

			// Only set if flag exists and wasn't set via command line
			if flag := pflag.Lookup(name); flag != nil && !flag.Changed {
				if err := flag.Value.Set(value); err != nil {
					log.Printf("Error setting flag %s: %v", name, err)
				}
			}
		}
	}

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

	// Split Comma-Separated list into an array
	allowedUsersArray := []string{}
	if len(*allowedUsers) > 0 {
		allowedUsersArray = strings.Split(*allowedUsers, ",")
	}

	log.Printf("Stakater Git WebHook Proxy started with provider '%s'\n", lowerProvider)

	if *insecureTLS {
		log.Println("WARNING: upstream TLS certificate verification is disabled (insecureSkipVerify=true)")
	}
	proxy.SetInsecureSkipVerify(*insecureTLS)

	p, err := proxy.NewProxy(*upstreamURL, allowedPathsArray, lowerProvider, *secret, ignoredUsersArray, allowedUsersArray)
	if err != nil {
		log.Fatal(err)
	}

	if err := p.Run(*listenAddress); err != nil {
		log.Fatal(err)
	}
}
