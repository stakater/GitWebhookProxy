package main

import (
	"log"
	//"os"
	"strings"

	config "github.com/rmenn/GitWebhookProxy/conf"
	"github.com/rmenn/GitWebhookProxy/pkg/proxy"
)

func validateRequiredFlags(upstreamURL, frontEndURL string) bool {
	isValid := true
	if len(strings.TrimSpace(upstreamURL)) == 0 {
		log.Println("Required flag 'upstreamURL' not specified")
		isValid = false
	}
	if len(strings.TrimSpace(frontEndURL)) == 0 {
		log.Println("Required Config 'frontEndURL' is not specified")
		isValid = false
	}
	return isValid
}

func main() {
	port, ps := config.Init()
	var Proxies []*proxy.Proxy
	for _, p := range ps {
		ok := validateRequiredFlags(p.UpstreamURL, p.FrontEndURL)
		if !ok {
			log.Fatal("required config not present")
		}
		pr, err := proxy.NewProxy(p.FrontEndURL, p.UpstreamURL, p.AllowedPaths, p.Provider,
			p.Secret, p.IgnoredUsers)
		if err != nil {
			log.Fatal("unable to configure proxy")
		}
		Proxies = append(Proxies, pr)
	}
	r := proxy.NewProxyRouter(Proxies)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
