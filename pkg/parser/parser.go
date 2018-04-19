package parser

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/stakater/JenkinsProxy/pkg/providers"
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

	if body, err := ioutil.ReadAll(req.Body); err != nil {
		return nil, err
	} else {
		hook.Payload = body
	}

	return hook, nil
}
