package main

import (
	"flag"
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
	upstreamUrl   = flag.String("upstreamUrl", "", "URL to which the proxy requests will be forwarded")
	secret        = flag.String("secret", "", "Secret of the Webhook API")
)

var allowedPaths arrayFlags

func main() {
	//TODO: paths allowed
	allowedPaths.Set("/")
	flag.Var(&allowedPaths, "allowedPaths", "Paths allowed to be forwarded via proxy")

	flag.Parse()
	log.Println("Stakater Git WebHook Proxy started ...")

	proxy.NewProxy(*listenAddress, *secret)
}
