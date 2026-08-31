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
	// GetHeaderKeys lists headers a request MUST carry. The parser rejects a
	// request missing any of them.
	GetHeaderKeys() []string

	// GetOptionalHeaderKeys lists headers the parser should forward WHEN
	// PRESENT and never require.
	//
	// The distinction did not exist and was needed the moment a provider
	// accepted more than one signature header. GitHub sends both
	// `X-Hub-Signature` and `X-Hub-Signature-256`, so listing the SHA-256 one
	// as required made BOTH mandatory — and a delivery carrying only the
	// modern header, which is what a receiver should prefer, was rejected
	// before `Validate` ever saw it. Caught by running the built image against
	// a real signed payload rather than by the unit tests, which call
	// `Validate` directly and never go through the parser.
	GetOptionalHeaderKeys() []string

	Validate(hook Hook) bool
	GetCommitter(hook Hook) string
	GetProviderName() string
}

func assertProviderImplementations() {
	var _ Provider = (*GithubProvider)(nil)
	var _ Provider = (*GitlabProvider)(nil)
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
