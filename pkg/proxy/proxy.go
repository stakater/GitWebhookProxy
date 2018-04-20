package proxy

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/stakater/JenkinsProxy/pkg/parser"
	"github.com/stakater/JenkinsProxy/pkg/providers"
)

type Proxy struct {
	provider     string
	allowedPaths []string
	secret       string
}

func (p *Proxy) isPathAllowed(path string) bool {
	// All paths allowed
	if len(p.allowedPaths) == 0 {
		return true
	}

	// Check if given passed exists in allowedPaths
	for _, p := range p.allowedPaths {
		if strings.TrimSuffix(p, "/") == strings.TrimSuffix(path, "/") {
			return true
		}
	}
	return false
}

func (p *Proxy) proxyRequest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	if !p.isPathAllowed(r.URL.Path) {
		log.Printf("Not allowed to proxy path: '%s'", r.URL.Path)
		http.Error(w, "Not allowed to proxy path: '"+r.URL.Path+"'", http.StatusForbidden)
		return
	}

	provider, err := providers.NewProvider(p.provider, p.secret)
	if err != nil {
		log.Printf("Error creating provider: %s", err)
		http.Error(w, "Error creating Provider", http.StatusInternalServerError)
		return
	}

	hook, err := parser.Parse(r, provider)
	if err != nil {
		log.Printf("Eror Parsing Hook: %s", err)
		http.Error(w, "Error parsing Hook: "+err.Error(), http.StatusBadRequest)
		return
	}

	if !provider.Validate(*hook) {
		log.Printf("Eror Validating Hook: %s", err)
		http.Error(w, "Error validating Hook", http.StatusBadRequest)
		return
	}

	//TODO: Create new request and forward to upstream
	s := fmt.Sprintf("%v", hook)
	w.Write([]byte(s))
}

// Health Check Endpoint
func (p *Proxy) health(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
}

func NewProxy(listenAddress string, allowedPaths []string, provider string, secret string) {
	// Validate Params
	if len(strings.TrimSpace(listenAddress)) == 0 {
		panic("Cannot create Proxy with empty listenAddress")
	}
	if len(strings.TrimSpace(secret)) == 0 {
		panic("Cannot create Proxy with empty secret")
	}

	proxy := Proxy{
		provider:     provider,
		allowedPaths: allowedPaths,
		secret:       secret,
	}

	router := httprouter.New()
	router.GET("/health", proxy.health)
	router.POST("/*path", proxy.proxyRequest)

	log.Printf("Listening at: %s", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, router))
}
