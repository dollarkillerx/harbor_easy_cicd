package models

type GitTask struct {
	Model
	GitAddress string       `json:"git_address"`
	GitType    string       `json:"git_type"`
	Repository string       `json:"repository"`
	Branch     string       `json:"branch"`   // 为空 则匹配所有
	Tag        string       `json:"tag"`      // 更具tag 匹配
	Comment    string       `json:"comment"`  // 根据comment 匹配
	Matching   MatchingType `json:"matching"` // 匹配方式

	TaskName  string `json:"task_name"`
	Path      string `json:"path"`
	Cmd       string `json:"cmd"`
	Heartbeat string `json:"heartbeat"`

	Run         bool `json:"run"`           // 是否在运行
	LastRunTime int  `json:"last_run_time"` // 最后运行时间
}

type MatchingType string

const (
	MathTag     MatchingType = "MathTag"
	MathComment MatchingType = "MathComment"
)

type GitType string

const (
	Github GitType = "Github"
	Gitea  GitType = "Gitea"
	Gitlab GitType = "Gitlab"
)

type GitlabPush struct {
	ObjectKind string `json:"object_kind"` // push
	EventName  string `json:"event_name"`
	Ref        string `json:"ref"` // branch
	Commits    []struct {
		Message string `json:"message"`
	} `json:"commits"`
	Repository struct {
		Name string `json:"name"`
		Url  string `json:"url"` // path (user/project)
	} `json:"repository"`
}

type GitlabTag struct {
	ObjectKind   string `json:"object_kind"` // tag_push
	EventName    string `json:"event_name"`
	Ref          string `json:"ref"` // tag name
	RefProtected bool   `json:"ref_protected"`
	Message      string `json:"message"`
	Repository   struct {
		Name string `json:"name"`
		Url  string `json:"url"` // path (user/project)
	} `json:"repository"`
}

// Gitea 相同
type GitHubPush struct {
	Ref        string `json:"ref"` // refs/heads/main
	Repository struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`      // project name
		FullName string `json:"full_name"` // file project name
	} `json:"repository"`
	HeadCommit struct {
		Message string `json:"message"` // comment
	} `json:"head_commit"`
}

type GithubTag struct {
	Ref        string `json:"ref"`      // tag 0.0.1
	RefType    string `json:"ref_type"` // type: tag
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
}

type GiteaTag struct {
	Action  string `json:"action"` // published
	Release struct {
		TagName string `json:"tag_name"`
	} `json:"release"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
}
