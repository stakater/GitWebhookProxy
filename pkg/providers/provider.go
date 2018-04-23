package providers

import (
	"errors"
	"strings"
)

const (
	GithubProviderKind = "github"
	GitlabProviderKind = "gitlab"
)

type Provider interface {
	GetHeaderKeys() []string
	Validate(hook Hook) bool
}

func NewProvider(provider string, secret string) (Provider, error) {
	switch strings.ToLower(provider) {
	case GithubProviderKind:
		return NewGithubProvider(secret)
	case GitlabProviderKind:
		return NewGitlabProvider(secret)
	default:
		return nil, errors.New("Unknown Provider git '" + provider + "' specified")
	}
}

type Hook struct {
	Payload []byte
	Headers map[string]string
}
