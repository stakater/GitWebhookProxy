package providers

import (
	"errors"
	"strings"
)

const (
	GithubProviderKind = "github"
	GitlabProviderKind = "gitlab"
	ContentTypeHeader  = "Content-Type"
)

type Provider interface {
	GetHeaderKeys() []string
	Validate(hook Hook) bool
}

func NewProvider(provider string, secret string) (Provider, error) {
	if len(provider) == 0 {
		return nil, errors.New("Empty provider string specified")
	}

	switch strings.ToLower(provider) {
	case GithubProviderKind:
		return NewGithubProvider(secret)
	case GitlabProviderKind:
		return NewGitlabProvider(secret)
	default:
		return nil, errors.New("Unknown Git Provider '" + provider + "' specified")
	}
}

type Hook struct {
	Payload       []byte
	Headers       map[string]string
	RequestMethod string
}
