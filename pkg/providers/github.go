package providers

import (
	"errors"
	"strings"
)

const (
	XGitHubSignature = "X-GitHub-Signature"
	XGitHubEvent     = "X-GitHub-Event"
	XGitHubDelivery  = "X-GitHub-Delivery"
)

type GithubProvider struct {
	secret string
}

func NewGithubProvider(secret string) (*GithubProvider, error) {
	if len(strings.TrimSpace(secret)) == 0 {
		return nil, errors.New("Cannot create github provider with empty secret")
	}

	return &GithubProvider{
		secret: secret,
	}, nil
}

func (p *GithubProvider) GetHeaderKeys() []string {
	return []string{
		XGitHubDelivery,
		XGitHubEvent,
		XGitHubSignature,
	}
}

func (p *GithubProvider) GetTokenHeaderKey() string {
	return XGitHubSignature
}

func (p *GithubProvider) Validate(hook Hook) bool {
	//TODO: Validate
	return true
}
