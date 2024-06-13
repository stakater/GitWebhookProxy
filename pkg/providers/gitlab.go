package providers

import (
	"encoding/json"
	"log"
	"strings"
)

// Header constants
const (
	XGitlabToken = "X-Gitlab-Token"
	XGitlabEvent = "X-Gitlab-Event"
	GitlabName   = "gitlab"
)

const (
	GitlabPushEvent         Event = "Push Hook"
	GitlabMergeRequestEvent Event = "Merge Request Hook"
)

type GitlabProvider struct {
	secret string
}

func NewGitlabProvider(secret string) (*GitlabProvider, error) {
	return &GitlabProvider{
		secret: secret,
	}, nil
}

func (p *GitlabProvider) GetProviderName() string {
	return GitlabName
}

// Not adding XGitlabToken will make token validation optional
func (p *GitlabProvider) GetHeaderKeys() []string {
	if len(strings.TrimSpace(p.secret)) > 0 {
		return []string{
			XGitlabEvent,
			XGitlabToken,
			ContentTypeHeader,
		}
	}

	return []string{
		XGitlabEvent,
		ContentTypeHeader,
	}
}

// Gitlab token validation:
// https://docs.gitlab.com/ee/user/project/integrations/webhooks.html#secret-token
func (p *GitlabProvider) Validate(hook Hook) bool {
	token := hook.Headers[XGitlabToken]
	// Validation fails if secret is configured but did not receive from gitlab
	if len(token) <= 0 {
		return false
	}

	return strings.TrimSpace(token) == strings.TrimSpace(p.secret)
}

func (p *GitlabProvider) GetCommitter(hook Hook) string {
	var payloadData GitlabPushPayload
	if err := json.Unmarshal(hook.Payload, &payloadData); err != nil {
		log.Printf("Gitlab hook payload unmarshalling failed: %v", err)
		return ""
	}

	eventType := Event(hook.Headers[XGitlabEvent])
	switch eventType {
	case GitlabPushEvent:
		return payloadData.Username
	}
	return ""
}
