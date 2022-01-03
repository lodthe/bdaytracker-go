package mailing

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	limiter "github.com/chatex-com/rate-limiter"
	"github.com/chatex-com/rate-limiter/pkg/config"
	"github.com/lodthe/bdaytracker-go/internal/conf"
	"github.com/petuhovskiy/telegram"
	"github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/internal/tgstate"
	"github.com/lodthe/bdaytracker-go/internal/tgview"
	"github.com/lodthe/bdaytracker-go/internal/usersession"
)

type StateGetter interface {
	GetAll() ([]*tgstate.State, error)
}

type SessionIssuer interface {
	Issue(telegramID int, update *telegram.Update) (s *usersession.Session, release func(), err error)
}

type Service struct {
	cfg *conf.Mailing

	stateGetter StateGetter
	issuer      SessionIssuer

	rateLimiter *limiter.RateLimiter
}

func NewService(cfg *conf.Mailing, stateGetter StateGetter, issuer SessionIssuer) *Service {
	rCfg := config.NewConfigWithQuotas([]*config.Quota{
		config.NewQuota(cfg.MaxRemindersPerSecond, time.Second),
	})
	rCfg.Concurrency = 1

	rateLimiter, _ := limiter.NewRateLimiter(rCfg)
	rateLimiter.Start()

	return &Service{
		stateGetter: stateGetter,
		issuer:      issuer,
		rateLimiter: rateLimiter,
	}
}

func (s *Service) Run(ctx context.Context) {
	const window = time.Minute * 15
	ticker := time.NewTicker(window)
	defer ticker.Stop()

	logrus.Info("started the mailing service")

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			s.sendReminders()
		}
	}
}

func (s *Service) readStates() ([]*tgstate.State, error) {
	states, err := s.stateGetter.GetAll()
	if err != nil {
		logrus.WithError(err).Error("failed to get all states")
		return nil, err
	}

	const day = time.Hour * 24
	today := time.Now().UTC().Truncate(day)

	var filtered []*tgstate.State
	for _, st := range states {
		if st.CannotReceiveMessages || st.LastNotificationAt.After(today) {
			continue
		}

		filtered = append(filtered, st)
	}

	return filtered, nil
}

func (s *Service) sendReminders() {
	if time.Now().Hour() < s.cfg.StartHour {
		return
	}
	if time.Now().Hour() > s.cfg.EndHour {
		return
	}

	states, err := s.readStates()
	if err != nil {
		logrus.WithError(err).Error("failed to read states")
		return
	}

	logrus.Debug("read states to send reminders")

	var sentReminders uint64
	wg := sync.WaitGroup{}

	for i := range states {
		wg.Add(1)

		go func(telegramID int) {
			defer wg.Done()

			logger := logrus.WithField("telegram_id", telegramID)

			session, release, err := s.issuer.Issue(telegramID, nil)
			if err != nil {
				logger.WithError(err).Error("failed to create a new session")
				return
			}
			defer release()

			response := <-s.rateLimiter.Execute(func() (interface{}, error) {
				bdayCount, err := tgview.Reminders{}.WishYourFriendsHappyBirthday(session)

				if bdayCount > 0 && err != nil {
					atomic.AddUint64(&sentReminders, 1)
				}

				return nil, err
			})

			if response.Error != nil {
				logger.WithError(response.Error).Error("failed to send a reminder to the user")
			} else {
				session.State.LastNotificationAt = time.Now()
			}

			err = session.SaveState()
			if err != nil {
				logger.WithError(err).Error("failed to save state")
			}
		}(states[i].TelegramID)
	}

	wg.Wait()

	logrus.WithField("reminder_count", sentReminders).Info("sent birthday reminders")
}
