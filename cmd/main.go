package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/dollarkillerx/harbor_easy_cicd/internal/conf"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/sdk/config"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/sdk/logs"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/server"
)

var configFilename string
var configDirs string

func init() {
	const (
		defaultConfigFilename = "config"
		configUsage           = "Name of the configs file, without extension"
		defaultConfigDirs     = "./,./configs/"
		configDirUsage        = "Directories to search for configs file, separated by ','"
	)
	flag.StringVar(&configFilename, "c", defaultConfigFilename, configUsage)
	flag.StringVar(&configFilename, "config", defaultConfigFilename, configUsage)
	flag.StringVar(&configDirs, "cPath", defaultConfigDirs, configDirUsage)
}

func main() {
	flag.Parse()

	// config
	var appConfig conf.Config
	err := config.InitConfiguration(configFilename, strings.Split(configDirs, ","), &appConfig)
	if err != nil {
		panic(err)
	}
	indent, err := json.MarshalIndent(appConfig, "", "  ")
	if err == nil {
		fmt.Println(string(indent))
	}
	fmt.Println("Config loaded successfully!")
	// 基础依赖初始化
	// 初始化日志
	logs.InitLog(appConfig.LoggerConfig)

	ser := server.NewServer(&appConfig)
	if err := ser.Run(); err != nil {
		panic(err)
	}
}
