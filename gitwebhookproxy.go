package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/stakater/JenkinsProxy/pkg/proxy"
)

// For allowing multiple flag values
type arrayFlags []string

func (af *arrayFlags) String() string {
	return strings.Join(*af, " ")
}

func (af *arrayFlags) Set(value string) error {
	*af = append(*af, strings.TrimSpace(value))
	return nil
}

// end

// proxy --listenAddress 8008 --provider github --secret mysecret
var (
	listenAddress = flag.String("listen", ":8080", "Address on which the proxy listens.")
	upstreamUrl   = flag.String("upstreamUrl", "", "URL to which the proxy requests will be forwarded (required)")
	secret        = flag.String("secret", "", "Secret of the Webhook API (required)")
	provider      = flag.String("provider", "github", "Git Provider which generates the Webhook")
)

var allowedPaths arrayFlags

func validateRequiredFlags() {
	isValid := true
	if len(strings.TrimSpace(*upstreamUrl)) == 0 {
		fmt.Println("Required flag 'upstreamUrl' not specified")
		isValid = false
	}
	if len(strings.TrimSpace(*secret)) == 0 {
		fmt.Println("Required flag 'secret' not specified")
		isValid = false
	}

	if !isValid {
		fmt.Println("")
		flag.Usage()
		fmt.Println("")

		panic("See Flag Usage")
	}
}

func main() {
	flag.Var(&allowedPaths, "allowPath", "Paths allowed to be forwarded via proxy. (All paths are allowed by default)")

	flag.Parse()
	validateRequiredFlags()
	lowerProvider := strings.ToLower(*provider)

	log.Println("Stakater Git WebHook Proxy started ...")
	proxy.NewProxy(*listenAddress, allowedPaths, lowerProvider, *secret)
}
