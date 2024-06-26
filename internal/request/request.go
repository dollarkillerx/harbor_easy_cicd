package request

import (
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
