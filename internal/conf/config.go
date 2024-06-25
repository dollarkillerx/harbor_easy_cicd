package conf

import "github.com/dollarkillerx/harbor_easy_cicd/internal/sdk/config"

type EasyCiCDConfig struct {
	AuthToken     string // auth token
	Address       string // address
	HarborAddress string

	LoggerConfig   config.LoggerConfig
	TelegramConfig *TelegramConfig // configs
	Tasks          []Task          // task任务
}

type Task struct {
	HarborKey string
	TaskName  string
	Path      string
	Cmd       string
	Heartbeat string
}

type TelegramConfig struct {
	BoltDBPath        string // bolt
	BotToken          string // bot token
	UserRegisterToken string // 用户注册token
}
