package handle

import (
	"github.com/petuhovskiy/telegram"

	"github.com/lodthe/bdaytracker-go/tg"
)

func StartHandling(general tg.General, updates <-chan telegram.Update) {
	for update := range updates {
		switch {
		case update.Message != nil:
			go handleMessage(general, update.Message)

		case update.CallbackQuery != nil:
			go handleCallback(general, update.CallbackQuery)
		}
	}
}

func handleMessage(general tg.General, msg *telegram.Message) {

}

func handleCallback(general tg.General, clb *telegram.CallbackQuery) {

}
