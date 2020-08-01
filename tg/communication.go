package tg

import (
	"strconv"

	"github.com/petuhovskiy/telegram"
	log "github.com/sirupsen/logrus"
)

const parseMode = "HTML"

func (s *Session) sendMessage(text string, keyboard telegram.AnyKeyboard) error {
	_, err := s.Bot.SendMessage(&telegram.SendMessageRequest{
		ChatID:                strconv.Itoa(s.TelegramID),
		Text:                  text,
		ParseMode:             parseMode,
		DisableWebPagePreview: true,
		ReplyMarkup:           keyboard,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"telegram_id":  s.TelegramID,
			"message_text": text,
		}).WithError(err).Error("failed to send the message")
	}
	return err
}

func (s *Session) SendText(text string) error {
	return s.sendMessage(text, nil)
}
