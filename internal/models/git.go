package models

type GitTask struct {
	Model
	Branch   string       // 为空 则匹配所有
	Tag      string       // 更具tag 匹配
	Comment  string       // 根据comment 匹配
	Matching MatchingType // 匹配方式

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

// github gitee
type GitPush struct {
	Ref        string `json:"ref"` // refs/heads/main
	HeadCommit struct {
		Id      string `json:"id"`
		Message string `json:"message"`
	} `json:"head_commit"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"` // path (user/project)
	} `json:"repository"`
}

type GithubTag struct {
	Ref          string  `json:"ref"`           // tag name
	RefType      string  `json:"ref_type"`      // type: tag
	MasterBranch string  `json:"master_branch"` // branch
	Description  *string `json:"description"`
	Repository   struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"` // path (user/project)
	} `json:"repository"`
}

type GiteeTag struct {
	Action  string `json:"action"`
	Release struct {
		TagName         string `json:"tag_name"`         // tag name
		TargetCommitish string `json:"target_commitish"` // 分支
		Name            string `json:"name"`
		Body            string `json:"body"`
	} `json:"release"`
	Repository struct {
		Name        string `json:"name"`
		FullName    string `json:"full_name"` // path (user/project)
		Description string `json:"description"`
	} `json:"repository"`
}

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
