package proxy

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jbcom/GitWebhookProxy/pkg/parser"
	"github.com/jbcom/GitWebhookProxy/pkg/providers"
	"github.com/jbcom/GitWebhookProxy/pkg/utils"
	"github.com/julienschmidt/httprouter"
)

var (
	transport = &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient = &http.Client{
		Timeout:   time.Second * 30,
		Transport: transport,
	}
)

type Proxy struct {
	provider     string
	upstreamURL  string
	allowedPaths []string
	secret       string
	ignoredUsers []string
	allowedUsers []string
}

func (p *Proxy) isPathAllowed(path string) bool {
	// All paths allowed
	if len(p.allowedPaths) == 0 {
		return true
	}

	// Check if given passed exists in allowedPaths
	for _, p := range p.allowedPaths {
		allowedPath := strings.TrimSpace(p)
		incomingPath := strings.TrimSpace(path)
		if strings.TrimSuffix(allowedPath, "/") ==
			strings.TrimSuffix(incomingPath, "/") || strings.HasPrefix(incomingPath, allowedPath) {
			return true
		}
	}
	return false
}

func (p *Proxy) isIgnoredUser(committer string) bool {
	if len(p.ignoredUsers) > 0 {
		if exists, _ := utils.InArray(p.ignoredUsers, committer); exists {
			return true
		}
	}

	if committer == "" && p.provider == providers.GithubName {
		return true
	}

	return false
}

func (p *Proxy) isAllowedUser(committer string) bool {
	if len(p.allowedUsers) > 0 {
		if exists, _ := utils.InArray(p.allowedUsers, committer); exists {
			return true
		}

		return false
	}

	return true
}

func (p *Proxy) redirect(hook *providers.Hook, redirectURL string) (*http.Response, error) {
	if hook == nil {
		return nil, errors.New("Cannot redirect with nil Hook")
	}

	// Parse url to check validity
	url, err := url.Parse(redirectURL)
	if err != nil {
		return nil, err
	}

	// Assign default scheme as http if not specified
	if url.Scheme == "" {
		url.Scheme = "http"
	}

	// Create Redirect request
	req, err := http.NewRequest(hook.RequestMethod, url.String(), bytes.NewBuffer(hook.Payload))

	if err != nil {
		return nil, err
	}

	// Set Headers from hook
	for key, value := range hook.Headers {
		req.Header.Add(key, value)
	}

	return httpClient.Do(req)

}

func (p *Proxy) proxyRequest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	redirectURL := p.upstreamURL + r.URL.Path

	if r.URL.RawQuery != "" {
		redirectURL += "?" + r.URL.RawQuery
	}

	log.Printf("Proxying Request from '%s', to upstream '%s'\n", r.URL, redirectURL)

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
		log.Printf("Error Parsing Hook: %s", err)
		http.Error(w, "Error parsing Hook: "+err.Error(), http.StatusBadRequest)
		return
	}

	// SIGNATURE FIRST. AUTHORISATION SECOND. THIS ORDER WAS THE OTHER WAY AROUND
	// AND THAT WAS A HOLE.
	//
	// The user filter used to run first and answer 200 to anything it decided
	// to ignore — before the signature was ever checked. So an UNSIGNED request
	// whose payload carried no committer, or a committer on the ignore list,
	// got a 200 from a proxy that had not authenticated it at all.
	//
	// Measured against the built image, `GWP_SECRET` set, no signature header
	// whatsoever:
	//
	//     Incoming request from user:
	//     Ignoring request for user:
	//     -> HTTP 200
	//
	// It never reached upstream, so it was not a forwarding bug. It was worse
	// in a quieter way: a public endpoint reporting success for a request it
	// had refused to authenticate, which is exactly what a probe looks for.
	//
	// The committer name is read from the PAYLOAD, and the payload is only
	// trustworthy after the signature says so. Deciding anything on it first —
	// including deciding to ignore it — is deciding on attacker-supplied data.
	if len(strings.TrimSpace(p.secret)) > 0 && !provider.Validate(*hook) {
		log.Printf("Error Validating Hook for '%s'", r.URL)
		http.Error(w, "Error validating Hook", http.StatusBadRequest)
		return
	}

	committer := provider.GetCommitter(*hook)
	log.Printf("Incoming request from user: %s", committer)
	if p.isIgnoredUser(committer) || (!p.isAllowedUser(committer)) {
		log.Printf("Ignoring request for user: %s", committer)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Ignoring request for user: %s", committer)))
		return
	}

	resp, errs := p.redirect(hook, redirectURL)
	if errs != nil {
		log.Printf("Error Redirecting '%s' to upstream '%s': %s\n", r.URL, redirectURL, errs)
		http.Error(w, "Error Redirecting '"+r.URL.String()+"' to upstream '"+redirectURL+"'", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode >= 400 {
		log.Printf("Error Redirecting '%s' to upstream '%s', Upstream Redirect Status: %s\n", r.URL, redirectURL, resp.Status)
		http.Error(w, "Error Redirecting '"+r.URL.String()+"' to upstream '"+redirectURL+"' Upstream Redirect Status:"+resp.Status, resp.StatusCode)
		return
	}

	log.Printf("Redirected incomming request '%s' to '%s' with Response: '%s'\n",
		r.URL, redirectURL, resp.Status)

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error Reading upstream '%s' response body\n", r.URL)
		http.Error(w, "Error Reading upstream '"+redirectURL+"' Response body", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
}

// Health Check Endpoint
func (p *Proxy) health(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("I'm Healthy and I know it! ;) "))
}

// Run starts Proxy server
func (p *Proxy) Run(listenAddress string) error {
	if len(strings.TrimSpace(listenAddress)) == 0 {
		panic("Cannot create Proxy with empty listenAddress")
	}

	router := httprouter.New()
	router.GET("/health", p.health)
	router.POST("/*path", p.proxyRequest)

	log.Printf("Listening at: %s", listenAddress)
	return http.ListenAndServe(listenAddress, router)
}

func NewProxy(upstreamURL string, allowedPaths []string,
	provider string, secret string, ignoredUsers []string) (*Proxy, error) {
	// Validate Params
	if len(strings.TrimSpace(upstreamURL)) == 0 {
		return nil, errors.New("Cannot create Proxy with empty upstreamURL")
	}
	if len(strings.TrimSpace(provider)) == 0 {
		return nil, errors.New("Cannot create Proxy with empty provider")
	}
	if allowedPaths == nil {
		return nil, errors.New("Cannot create Proxy with nil allowedPaths")
	}

	return &Proxy{
		provider:     provider,
		upstreamURL:  upstreamURL,
		allowedPaths: allowedPaths,
		secret:       secret,
		ignoredUsers: ignoredUsers,
	}, nil
}
