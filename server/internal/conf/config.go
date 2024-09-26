package conf

import "github.com/dollarkillerx/harbor_easy_cicd/internal/sdk/config"

type Config struct {
	AuthToken     string // auth token
	Address       string // address
	HarborAddress string

	AdminAuth                 AdminAuth
	PostgresConfig            config.PostgresConfiguration
	LoggerConfig              config.LoggerConfig
	TelegramUserRegisterToken string // 用户注册token
}

type AdminAuth struct {
	Token string
}
