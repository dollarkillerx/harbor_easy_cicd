package server

import (
	"github.com/dollarkillerx/harbor_easy_cicd/internal/conf"
	"github.com/dollarkillerx/harbor_easy_cicd/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	conf *conf.Config
	app  *gin.Engine

	sendData chan string
	task     []Task
}

func NewServer(conf *conf.Config) *Server {
	return &Server{
		conf:     conf,
		sendData: make(chan string, 10),
		task:     make([]Task, 0),
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
}
