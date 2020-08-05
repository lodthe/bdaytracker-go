package tg

import (
	"errors"
	"strconv"

	"github.com/petuhovskiy/telegram"
	"github.com/petuhovskiy/telegram/markup"
	log "github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/static"
)

const parseMode = "HTML"

func (s *Session) AnswerOnLastCallback() {
	if s.LastUpdate == nil || s.LastUpdate.CallbackQuery == nil {
		return
	}
	s.Bot.AnswerCallbackQuery(&telegram.AnswerCallbackQueryRequest{
		CallbackQueryID: s.LastUpdate.CallbackQuery.ID,
	})
}

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

func (s *Session) editInlineMessage(text string, keyboard *telegram.InlineKeyboardMarkup) error {
	_, err := s.Bot.EditMessageText(&telegram.EditMessageTextRequest{
		ChatID:                strconv.Itoa(s.TelegramID),
		MessageID:             s.LastUpdate.CallbackQuery.Message.MessageID,
		InlineMessageID:       s.LastUpdate.CallbackQuery.InlineMessageID,
		Text:                  text,
		ParseMode:             parseMode,
		DisableWebPagePreview: true,
		ReplyMarkup:           keyboard,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"telegram_id":    s.TelegramID,
			"message_text":   text,
			"callback_query": s.LastUpdate.CallbackQuery,
		}).WithError(err).Error("failed to edit the message")
	}
	return err
}

func (s *Session) SendText(text string, keyboard ...telegram.AnyKeyboard) error {
	if len(keyboard) == 0 {
		return s.sendMessage(text, nil)
	}

	switch buttons := keyboard[0].(type) {
	case [][]telegram.InlineKeyboardButton:
		return s.sendMessage(text, markup.InlineKeyboard(buttons))

	case [][]telegram.KeyboardButton:
		return s.sendMessage(text, telegram.ReplyKeyboardMarkup{
			Keyboard:        buttons,
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
			Selective:       true,
		})

	default:
		err := errors.New("unknown keyboard type")
		log.WithField("keyboard", keyboard).WithError(err).Error("failed to send a telegram message")
		return err
	}
}

func (s *Session) SendEditText(text string, keyboard telegram.AnyKeyboard, edit bool) error {
	if !edit || s.LastUpdate.CallbackQuery == nil {
		return s.SendText(text, keyboard)
	}

	switch buttons := keyboard.(type) {
	case [][]telegram.InlineKeyboardButton:
		return s.editInlineMessage(text, markup.InlineKeyboardMarkup(buttons))

	case [][]telegram.KeyboardButton:
		log.WithFields(log.Fields{
			"text":     text,
			"keyboard": keyboard,
		}).Error("trying to edit the message with not an inline markup")
		return s.SendText(text, keyboard)

	default:
		err := errors.New("unknown keyboard type")
		log.WithField("keyboard", keyboard).WithError(err).Error("failed to send a telegram message")
		return err
	}
}

func (s *Session) SendInlinePhoto(text string, file string, keyboard telegram.AnyKeyboard) error {
	_, err := s.Bot.SendPhoto(&telegram.SendPhotoRequest{
		ChatID:      strconv.Itoa(s.TelegramID),
		Photo:       static.NewFileReader(file),
		Caption:     text,
		ParseMode:   parseMode,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"telegram_id":  s.TelegramID,
			"message_text": text,
			"file":         file,
		}).WithError(err).Error("failed to send the message")
	}
	return err
}
