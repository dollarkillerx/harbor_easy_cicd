package client

import (
	"fmt"

	"github.com/dollarkillerx/harbor_easy_cicd/internal/sdk/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func PostgresClient(conf config.PostgresConfiguration, gormConfig *gorm.Config) (*gorm.DB, error) {
	if conf.TimeZone == "" {
		conf.TimeZone = "Asia/Shanghai"
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=%s", conf.Host, conf.User, conf.Password, conf.DBName, conf.Port, conf.TimeZone)
	if !conf.SSLMode {
		dsn += " sslmode=disable"
	}

	if gormConfig == nil {
		gormConfig = &gorm.Config{}
	}

	return gorm.Open(postgres.Open(dsn), gormConfig)
}
