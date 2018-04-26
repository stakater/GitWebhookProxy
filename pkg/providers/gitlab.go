package providers

import (
	"errors"
	"strings"
)

// Header constants
const (
	XGitlabToken = "X-Gitlab-Token"
	XGitlabEvent = "X-Gitlab-Event"
)

type GitlabProvider struct {
	secret string
}

func NewGitlabProvider(secret string) (*GitlabProvider, error) {
	if len(strings.TrimSpace(secret)) == 0 {
		return nil, errors.New("Cannot create github provider with empty secret")
	}

	return &GitlabProvider{
		secret: secret,
	}, nil
}

func (p *GitlabProvider) GetHeaderKeys() []string {
	return []string{
		XGitlabToken,
		XGitlabEvent,
	}
}

// Gitlab token validation:
// https://docs.gitlab.com/ee/user/project/integrations/webhooks.html#secret-token
func (p *GitlabProvider) Validate(hook Hook) bool {
	token := hook.Headers[XGitlabToken]
	if len(token) <= 0 {
		return false
	}

	return strings.TrimSpace(token) == strings.TrimSpace(p.secret)
}
