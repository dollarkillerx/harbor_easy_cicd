package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func (s *Server) noticeLog(key string, name string, message string) {
	log.Info().Msgf("Key: %s Task: %s Time: %s message: %s", key, name, time.Now().Format("2006-01-02 15:04:05"), message)
	if strings.TrimSpace(s.conf.TelegramUserRegisterToken) != "" {
		s.sendData <- fmt.Sprintf("Key: %s Task: %s Time: %s message: %s", key, name, time.Now().Format("2006-01-02 15:04:05"), message)
	}
}
