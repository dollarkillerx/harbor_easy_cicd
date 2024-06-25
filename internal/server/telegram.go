package server

import (
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"

	"strings"
)

func (s *Server) telegram() {
	if s.conf.TelegramConfig == nil {
		return
	}

	bot, err := tgbotapi.NewBotAPI(s.conf.TelegramConfig.BotToken)
	if err != nil {
		panic(err)
	}

	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	db, err := bolt.Open(s.conf.TelegramConfig.BoltDBPath, 0600, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Users"))
		return err
	})

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	go func() {
		for {
			select {
			case messageText := <-s.sendData:
				db.View(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte("Users"))
					b.ForEach(func(k, v []byte) error {
						chatID := btoi(k)
						msg := tgbotapi.NewMessage(chatID, messageText)
						bot.Send(msg)
						return nil
					})
					return nil
				})
			}
		}
	}()

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			update.Message.Text = strings.TrimSpace(update.Message.Text)
			if update.Message.Text != s.conf.TelegramConfig.UserRegisterToken {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "身份验证失败")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
				continue
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "身份验证成功")
			msg.ReplyToMessageID = update.Message.MessageID

			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("Users"))
				id := itob(update.Message.Chat.ID)
				return b.Put(id, []byte("active"))
			})

			bot.Send(msg)
		}
	}
}

func itob(v int64) []byte {
	b := make([]byte, 8)
	for i := uint(0); i < 8; i++ {
		b[i] = byte(v >> (i * 8))
	}
	return b
}

func btoi(b []byte) int64 {
	var v int64
	for i := uint(0); i < 8; i++ {
		v |= int64(b[i]) << (i * 8)
	}
	return v
}
