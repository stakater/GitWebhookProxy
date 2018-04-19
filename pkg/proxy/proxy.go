package proxy

import (
	"log"
	"net/http"

	"github.com/stakater/JenkinsProxy/pkg/providers"
)

type Proxy struct {
	secret string
}

func (p *Proxy) proxyRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Header.Get("X-GitHub-Delivery")))

	provider, err := providers.NewProvider("github", p.secret)
	if err != nil {
		log.Printf("ERORRRRRR %s", err)
		return
	}

	hook, err := provider.GetHook(r)
	if err != nil {
		log.Printf("ERORRRRRR %s", err)
		return
	}

	w.Write([]byte(hook.Id))
}

func NewProxy(listenAddress string, secret string) {

	proxy := Proxy{secret: secret}
	http.HandleFunc("/", proxy.proxyRequest)

	log.Printf("Listening at: %s", listenAddress)
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		panic(err)
	}
}
