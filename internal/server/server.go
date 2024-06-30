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

	postgresClient.AutoMigrate(&models.Task{}, &models.TaskLogs{}, &models.GitTask{})

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
	// 提供静态文件
	s.app.Static("/assets", "./dist/assets")

	// 提供单页应用的入口文件
	//s.app.LoadHTMLFiles("dist/index.html")

	s.app.POST("/hook", middleware.Auth(s.conf.AuthToken), s.webHook)
	s.app.POST("/hook_github", s.webHookGit)
	s.app.POST("/hook_gitee", s.webHookGit)
	s.app.POST("/hook_gitlab", s.webHookGit)
	//s.app.POST("/hook_git", middleware.Auth(s.conf.AuthToken), s.webHookGit)
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

	// 所有其他请求都返回 index.html
	// 捕获所有其他未匹配的路由，并返回 index.html
	//s.app.NoRoute(func(c *gin.Context) {
	//	path := c.Request.URL.Path
	//	// 如果请求路径以 /assets/ 开头，则直接返回 404 错误
	//	log.Info().Msgf("%v", path)
	//	// 否则返回 index.html
	//	c.HTML(http.StatusOK, "index.html", nil)
	//})
}
