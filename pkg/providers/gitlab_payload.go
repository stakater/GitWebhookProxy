package providers

// PushPayload contains the information for Gitlab's push hook event
type GitlabPushPayload struct {
	Kind          string `json:"push"`
	Before        string `json:"before"`
	After         string `json:"after"`
	Ref           string `json:"ref"`
	CheckoutSha   string `json:"checkout_sha"`
	UserId        int64  `json:"user_id"`
	Name          string `json:"user_name"`
	Username      string `json:"user_username"`
	Email         string `json:"user_email"`
	UserAvatarUrl string `json:"user_avatar"`
	ProjectId     int64  `json:"project_id"`
	Project       struct {
		ProjectId          int64   `json:"id"`
		ProjectName        string  `json:"name"`
		ProjectDescription string  `json:"description"`
		ProjectWebUrl      string  `json:"web_url"`
		ProjectAvatarUrl   *string `json:"avatar_url"`
		GitSshUrl          string  `json:"git_ssh_url"`
		GitHttpUrl         string  `json:"git_http_url"`
		Namespace          string  `json:"namespace"`
		VisibilityLevel    int64   `json:"visibility_level"`
		NamespacePath      string  `json:"path_with_namespace"`
		DefaultBranch      string  `json:"default_branch"`
		HomePageUrl        string  `json:"homepage"`
		ProjectUrl         string  `json:"url"`
		ProjectSshUrl      string  `json:"ssh_url"`
		ProjectHttpUrl     string  `json:"http_url"`
	} `json:"project"`
	Repository struct {
		RepoName            string `json:"name"`
		RepoUrl             string `json:"url"`
		RepoDescription     string `json:"description"`
		RepoHompageUrl      string `json:"homepage"`
		RepoHttpUrl         string `json:"git_http_url"`
		RepoSshUrl          string `json:"git_ssh_url"`
		RepoVisibilityLevel int64  `json:"visibility_level"`
	} `json:"repository"`
	Commits []struct {
		CommitId        string `json:"Commits"`
		CommitMessage   string `json:"message"`
		CommitTimestamp string `json:"timestamp"`
		CommitUrl       string `json:"url"`
		Author          struct {
			AutherName  string `json:"name"`
			AutherEmail string `json:"email"`
		} `json:"author"`
		CommitAdded    []string `json:"added"`
		CommitModified []string `json:"modified"`
		CommitRemoved  []string `json:"removed"`
	} `json:"commits"`
	TotalCommitsCount int64 `json:"total_commits_count"`
}
