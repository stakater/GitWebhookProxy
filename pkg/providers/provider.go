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
	GetEventType(hook Hook) Event
	IsCommitterCheckEvent(event Event) bool
	GetCommitter(hook Hook, eventType Event) string
	GetProviderName() string
}

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
		return nil, errors.New("unknown Git Provider '" + provider + "' specified")
	}
}

type Hook struct {
	Payload       []byte
	Headers       map[string]string
	RequestMethod string
}
