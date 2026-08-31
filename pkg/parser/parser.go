package parser

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/jbcom/GitWebhookProxy/pkg/providers"
)

func Parse(req *http.Request, provider providers.Provider) (*providers.Hook, error) {
	hook := &providers.Hook{
		Headers: make(map[string]string),
	}

	for _, header := range provider.GetHeaderKeys() {
		if req.Header.Get(header) != "" {
			hook.Headers[header] = req.Header.Get(header)
			continue
		}
		return nil, errors.New("Required header '" + header + "' not found in Request")
	}

	// OPTIONAL HEADERS ARE TAKEN WHEN PRESENT AND NEVER DEMANDED. A provider
	// that accepts more than one signature header — GitHub sends both
	// `X-Hub-Signature` and `X-Hub-Signature-256` — cannot express that through
	// the required list, where every entry is mandatory. Listing either as
	// required rejects a legitimate sender: SHA-256 excludes GitLab and older
	// GitHub Enterprise, SHA-1 excludes anyone who has dropped the legacy
	// header. The rule is "at least one, strongest wins", and only `Validate`
	// can state it.
	for _, header := range provider.GetOptionalHeaderKeys() {
		if value := req.Header.Get(header); value != "" {
			hook.Headers[header] = value
		}
	}

	if body, err := ioutil.ReadAll(req.Body); err != nil {
		return nil, err
	} else {
		hook.Payload = body
	}

	hook.RequestMethod = req.Method

	return hook, nil
}
