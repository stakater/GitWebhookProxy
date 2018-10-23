package providers

import (
	"encoding/json"
	"strings"

	"github.com/sirupsen/logrus"
)

// Header constants
const (
	XGitlabToken = "X-Gitlab-Token"
	XGitlabEvent = "X-Gitlab-Event"
)

const (
	GitlabPushEvent Event = "Push Hook"
)

type GitlabProvider struct {
	secret string
}

func NewGitlabProvider(secret string) (*GitlabProvider, error) {
	return &GitlabProvider{
		secret: secret,
	}, nil
}

// Not adding XGitlabToken to make token validation optional
func (p *GitlabProvider) GetHeaderKeys() []string {
	return []string{
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
	var payloadData GitlabPushPayload
	if err := json.Unmarshal(hook.Payload, &payloadData); err != nil {
		logrus.Errorf("Gitlab hook payload unmarshling failed")
		return ""
	}

	eventType := Event(hook.Headers[XGitlabEvent])
	switch eventType {
	case GitlabPushEvent:
		return payloadData.Username
	}
	return ""
}
