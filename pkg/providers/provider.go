package providers

import (
	"errors"
	"strings"
)

const (
	GithubProviderKind            = "github"
	GitlabProviderKind            = "gitlab"
	ContentTypeHeader             = "Content-Type"
	DefaultContentTypeHeaderValue = "application/json"
)

// Event defines a provider hook event type
type Event string

type Provider interface {
	GetHeaderKeys() []string
	Validate(hook Hook) bool
	GetCommitter(hook Hook) string
	GetProviderName() string
}

// assertProviderImplementations is used for compile-time verification that
// provider types implement the Provider interface
//
//nolint:unused // This function is used for compile-time type checking
func assertProviderImplementations() {
	var _ Provider = (*GithubProvider)(nil)
	var _ Provider = (*GitlabProvider)(nil)
}

func NewProvider(provider string, secret string) (Provider, error) {
	if len(provider) == 0 {
		return nil, errors.New("empty provider string specified")
	}

	switch strings.ToLower(provider) {
	case GithubProviderKind:
		return NewGithubProvider(secret)
	case GitlabProviderKind:
		return NewGitlabProvider(secret)
	default:
		return nil, errors.New("unknown git provider '" + provider + "' specified")
	}
}

type Hook struct {
	Payload       []byte
	Headers       map[string]string
	RequestMethod string
}
