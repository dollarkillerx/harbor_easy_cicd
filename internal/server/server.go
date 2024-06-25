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
}

func NewServer(conf *conf.Config) *Server {
	return &Server{
		conf:     conf,
		sendData: make(chan string, 10),
	}
}

func (s *Server) Run() error {
	go s.telegram()

	s.app = gin.Default()
	gin.SetMode(gin.ReleaseMode)

	s.app.Use(middleware.Cors())
	s.app.Use(middleware.Auth(s.conf.AuthToken))

	s.router()

	return s.app.Run(s.conf.Address)
}

func (s *Server) router() {
	s.app.POST("/hook", s.webHook)
}
