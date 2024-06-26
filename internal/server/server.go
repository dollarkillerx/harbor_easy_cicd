package server

import (
	"github.com/dollarkillerx/harbor_easy_cicd/internal/conf"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/middleware"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/models"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/sdk/client"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	conf *conf.Config
	app  *gin.Engine
	db   *gorm.DB

	sendData chan string
}

func NewServer(conf *conf.Config) *Server {
	postgresClient, err := client.PostgresClient(conf.PostgresConfig, nil)
	if err != nil {
		panic(err)
	}

	postgresClient.AutoMigrate(&models.Task{}, &models.TaskLogs{})

	return &Server{
		conf:     conf,
		db:       postgresClient,
		sendData: make(chan string, 10),
	}
}

func (s *Server) Run() error {
	go s.telegram()

	s.app = gin.Default()
	gin.SetMode(gin.ReleaseMode)

	s.app.Use(middleware.Cors())

	s.router()

	return s.app.Run(s.conf.Address)
}

func (s *Server) router() {
	s.app.POST("/hook", middleware.Auth(s.conf.AuthToken), s.webHook)
	s.app.GET("/heartbeat", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "success",
		})
	})

	backstage := s.app.Group("/backstage/public")
	{
		backstage.POST("/login", s.login)
	}

	backstageInternal := s.app.Group("/backstage/internal", middleware.Auth(s.conf.AdminAuth.Token))
	{
		backstageInternal.GET("/tasks", s.tasks) // 查询
		backstageInternal.POST("/task", s.task)  // 增 删 改
		backstageInternal.GET("/logs", s.logs)   // 增 删 改
	}
}
