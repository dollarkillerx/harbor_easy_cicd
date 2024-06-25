package server

import (
	"github.com/rs/zerolog/log"
	"time"
)

func (s *Server) noticeLog(key string, name string, message string) {
	log.Info().Msgf("Key: %s Task: %s Time: %s message: %s", key, name, time.Now().Format("2006-01-02 15:04:05"), message)
	if s.conf.TelegramConfig != nil {
		s.conf.TelegramConfig.BotToken
	}
}
