package providers

import (
	"net/http"
	"strings"
)

const (
	GithubProviderKind = "github"
	GitlabProviderKind = "gitlab"
)

type Provider interface {
	GetHook(req *http.Request) (Hook, error)
	validate() bool
}

func NewProvider(provider string, secret string) (Provider, error) {
	switch strings.ToLower(provider) {
	case GithubProviderKind:
		return NewGithubProvider(secret)
	// case GitlabProviderKind:
	// 	return NewGitlabProvider(secret)
	default:
		return NewGithubProvider(secret)
	}
}

type Hook struct {
	Id    string
	Event string
	// Token or Signature
	Token   string
	Payload []byte
}
