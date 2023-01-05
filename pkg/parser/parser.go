package parser

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/stakater/GitWebhookProxy/pkg/providers"
)

func Parse(req *http.Request, provider providers.Provider) (*providers.Hook, error) {
	hook := &providers.Hook{
		Headers: make(map[string]string),
	}

	for _, header := range provider.GetHeaderKeys() {
		if req.Header.Get(header) == "" {
			return nil, errors.New("Required header '" + header + "' not found in Request")
		}

		// Store required headers in the expected casing
		hook.Headers[header] = req.Header.Get(header)
	}

	for header := range req.Header {
		// Store the rest of the headers in any casing
		hook.Headers[header] = req.Header.Get(header)
	}

	if body, err := ioutil.ReadAll(req.Body); err != nil {
		return nil, err
	} else {
		hook.Payload = body
	}

	hook.RequestMethod = req.Method

	return hook, nil
}
