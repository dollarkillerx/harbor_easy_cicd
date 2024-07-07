package request

import (
	"github.com/dollarkillerx/harbor_easy_cicd/internal/models"
	"github.com/pkg/errors"
	"strings"
)

type TaskPayload struct {
	Type      TaskType `json:"type"`
	TaskId    uint     `json:"task_id"`    // del,edit,stop,start
	HarborKey string   `json:"harbor_key"` // harbor_key
	TaskName  string   `json:"task_name"`  // taskName
	Path      string   `json:"path"`       // exec path
	Cmd       string   `json:"cmd"`        // 执行命令
	Heartbeat string   `json:"heartbeat"`  // 心跳地址
}

func (t *TaskPayload) Validate() error {
	if t.Type == "" {
		return errors.New("type error")
	}

	t.HarborKey = strings.TrimSpace(t.HarborKey)
	t.TaskName = strings.TrimSpace(t.TaskName)
	t.Path = strings.TrimSpace(t.Path)
	t.Cmd = strings.TrimSpace(t.Cmd)
	t.Heartbeat = strings.TrimSpace(t.Heartbeat)

	if t.Type == TaskAdd || t.Type == TaskEdit {
		if t.HarborKey == "" {
			return errors.New("Harbor 镜像名称不能为空")
		}

		if t.TaskName == "" {
			return errors.New("TaskName 不能为空")
		}

		if t.Path == "" {
			return errors.New("Path 不能为空")
		}

		if t.Cmd == "" {
			return errors.New("Cmd 不能为空")
		}
	}

	return nil
}

type TaskType string

const (
	TaskAdd   TaskType = "TaskAdd"
	TaskDel   TaskType = "TaskDel"
	TaskEdit  TaskType = "TaskEdit"
	TaskStop  TaskType = "TaskStop"
	TaskStart TaskType = "TaskStart"
)

type GitTaskPayload struct {
	Type TaskType `json:"type"`

	TaskId     uint                `json:"task_id"`     // del,edit,stop,start
	GitType    models.GitType      `json:"harbor_key"`  // harbor_key
	GitAddress string              `json:"git_address"` // git 地址
	Repository string              `json:"repository"`  // 仓库
	Branch     string              `json:"branch"`      // 分支
	Tag        string              `json:"tag"`         //  tag 匹配
	Comment    string              `json:"comment"`     // comment 匹配
	Matching   models.MatchingType `json:"matching"`    // 匹配方式

	TaskName  string `json:"task_name"` // taskName
	Path      string `json:"path"`      // exec path
	Cmd       string `json:"cmd"`       // 执行命令
	Heartbeat string `json:"heartbeat"` // 心跳地址

	Run bool `json:"run"` // 是否在运行
}

func (t *GitTaskPayload) Validate() error {
	t.GitAddress = strings.TrimSpace(t.GitAddress)
	t.Repository = strings.TrimSpace(t.Repository)
	t.Branch = strings.TrimSpace(t.Branch)
	t.Tag = strings.TrimSpace(t.Tag)
	t.Comment = strings.TrimSpace(t.Comment)
	t.TaskName = strings.TrimSpace(t.TaskName)
	t.Path = strings.TrimSpace(t.Path)
	t.Cmd = strings.TrimSpace(t.Cmd)
	t.Heartbeat = strings.TrimSpace(t.Heartbeat)

	if t.Type == TaskAdd || t.Type == TaskEdit {
		if t.GitAddress == "" {
			return errors.New("GitAddress 不能为空")
		}

		if t.Repository == "" {
			return errors.New("Repository 不能为空")
		}

		if t.Branch == "" {
			return errors.New("Branch 不能为空")
		}

		if t.Tag == "" {
			return errors.New("Tag 不能为空")
		}

		if t.Comment == "" {
			return errors.New("Comment 不能为空")
		}

		if t.Tag == "" {
			return errors.New("Tag 不能为空")
		}

		if t.TaskName == "" {
			return errors.New("TaskName 不能为空")
		}

		if t.Path == "" {
			return errors.New("Path 不能为空")
		}

		if t.Cmd == "" {
			return errors.New("Cmd 不能为空")
		}
	}

	return nil
}
