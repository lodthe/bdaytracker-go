package tgview

import (
	"time"

	"github.com/lodthe/bdaytracker-go/tg"
)

type Notifications struct {
}

func (n Notifications) WishYourFriendsHappyBirthday(s *tg.Session) (birthdaysNumber int, err error) {
	text := "Сегодня День Рождения у твоих друзей!\n\n"

	_, month, day := time.Now().Date()

	for _, friend := range s.State.Friends {
		if friend.BDay != day || friend.BMonth != int(month) {
			continue
		}

		text += formatFriend(friend) + "\n"
		birthdaysNumber++
	}

	if birthdaysNumber == 0 {
		return 0, nil
	}

	err = s.SendText(text)
	if err != nil {
		return 0, err
	}

	return birthdaysNumber, nil
}
