package providers

import (
	"crypto/subtle"
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
// GetOptionalHeaderKeys is empty: GitLab sends exactly one token header and it
// is required, so there is nothing conditional to forward.
func (p *GitlabProvider) GetOptionalHeaderKeys() []string {
	return nil
}

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
// Validate compares GitLab's X-Gitlab-Token against the configured secret.
//
// CONSTANT-TIME, because `==` on a secret is a timing oracle. Go's string
// comparison returns as soon as two bytes differ, so the time a rejection takes
// is a function of how many leading characters the attacker guessed correctly.
// Over enough requests — and a webhook endpoint is public and accepts as many
// as you send it — that recovers the token one byte at a time.
//
// GitLab sends the token verbatim rather than a signature, so unlike GitHub
// there is no digest to compare and this string IS the credential. That makes
// the comparison the whole of the check, and the only place to get it right.
//
// `subtle.ConstantTimeCompare` also returns 0 for differing lengths without
// examining contents, which is the correct behaviour here: a wrong-length token
// is wrong regardless.
func (p *GitlabProvider) Validate(hook Hook) bool {
	token := strings.TrimSpace(hook.Headers[XGitlabToken])
	// Validation fails if secret is configured but did not receive from gitlab
	if len(token) == 0 {
		return false
	}

	secret := strings.TrimSpace(p.secret)
	return subtle.ConstantTimeCompare([]byte(token), []byte(secret)) == 1
}

func (p *GitlabProvider) GetCommitter(hook Hook) string {
	var payloadData GitlabPushPayload
	if err := json.Unmarshal(hook.Payload, &payloadData); err != nil {
		log.Printf("Gitlab hook payload unmarshalling failed")
		return ""
	}

	eventType := Event(hook.Headers[XGitlabEvent])
	switch eventType {
	case GitlabPushEvent:
		return payloadData.Username
	}
	return ""
}
