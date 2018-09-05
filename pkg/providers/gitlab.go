package providers

import (
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
	return &GitlabProvider{
		secret: secret,
	}, nil
}

func (p *GitlabProvider) GetHeaderKeys() []string {
	return []string{
		XGitlabToken,
		XGitlabEvent,
		ContentTypeHeader,
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

func (p *GitlabProvider) GetCommitter(hook Hook) string {
	return ""
}
