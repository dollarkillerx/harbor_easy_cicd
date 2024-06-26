package models

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	Model
	HarborKey string `json:"harbor_key"`
	TaskName  string `json:"task_name"`
	Path      string `json:"path"`
	Cmd       string `json:"cmd"`
	Heartbeat string `json:"heartbeat"`

	Run         bool `json:"run"`           // 是否在运行
	LastRunTime int  `json:"last_run_time"` // 最后运行时间
}

type TaskLogs struct {
	Model
	TaskId   uint   `json:"task_id"`
	TaskName string `json:"task_name"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}

type Model struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
