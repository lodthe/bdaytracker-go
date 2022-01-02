package notifications

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	limiter "github.com/chatex-com/rate-limiter"
	"github.com/chatex-com/rate-limiter/pkg/config"
	"github.com/sirupsen/logrus"

	"github.com/lodthe/bdaytracker-go/tg"
	"github.com/lodthe/bdaytracker-go/tg/sessionstorage"
	"github.com/lodthe/bdaytracker-go/tg/state"
	"github.com/lodthe/bdaytracker-go/tg/tgview"
)

const notificationsStartHour = 7
const notificationsEndHour = 10
const maxNotificationsInSecond = 15

type Service struct {
	stateRepo state.Repository
	general   *tg.General
	storage   *sessionstorage.Storage

	rateLimiter *limiter.RateLimiter
}

func newRateLimiter() *limiter.RateLimiter {
	cfg := config.NewConfigWithQuotas([]*config.Quota{
		config.NewQuota(maxNotificationsInSecond, time.Second),
	})
	cfg.Concurrency = 1

	rateLimiter, _ := limiter.NewRateLimiter(cfg)
	rateLimiter.Start()

	return rateLimiter
}

func NewService(repo state.Repository, general *tg.General, storage *sessionstorage.Storage) *Service {
	return &Service{
		stateRepo:   repo,
		general:     general,
		storage:     storage,
		rateLimiter: newRateLimiter(),
	}
}

func (s *Service) Run(ctx context.Context) {
	const window = time.Minute * 15
	ticker := time.NewTicker(window)

	logrus.Info("started the notification service")

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			s.sendNotifications()
		}
	}
}

func (s *Service) readStates() ([]*state.State, error) {
	states, err := s.stateRepo.GetAll()
	if err != nil {
		logrus.WithError(err).Error("failed to get all states")
		return nil, err
	}

	const day = time.Hour * 24
	today := time.Now().UTC().Truncate(day)

	var filtered []*state.State
	for _, s := range states {
		if s.CannotReceiveMessages || s.LastNotificationAt.After(today) {
			continue
		}

		filtered = append(filtered, s)
	}

	return filtered, nil
}

func (s *Service) sendNotifications() {
	if time.Now().Hour() < notificationsStartHour {
		return
	}
	if time.Now().Hour() > notificationsEndHour {
		return
	}

	states, err := s.readStates()
	if err != nil {
		logrus.WithError(err).Error("failed to read states")
		return
	}

	logrus.Info("read states to send notifications")

	var sentNotifications uint64
	wg := sync.WaitGroup{}

	for i := range states {
		wg.Add(1)

		go func(userTelegramID int) {
			lock := s.storage.AcquireLock(userTelegramID)
			lock.Lock()
			defer lock.Unlock()

			logger := logrus.WithField("telegram_id", userTelegramID)

			session, err := tg.NewSession(s.general.VKCli, s.general.Bot, s.general.Executor, s.stateRepo, userTelegramID, nil)
			if err != nil {
				logger.WithError(err).Error("failed to create a new session")
				wg.Done()
				return
			}

			response := <-s.rateLimiter.Execute(func() (interface{}, error) {
				birthdaysNumber, err := tgview.Notifications{}.WishYourFriendsHappyBirthday(session)

				if birthdaysNumber > 0 && err != nil {
					atomic.AddUint64(&sentNotifications, 1)
				}

				return nil, err
			})

			if response.Error != nil {
				logger.WithError(response.Error).Error("failed to send a notification to the user")
			} else {
				session.State.LastNotificationAt = time.Now()
			}

			s.stateRepo.Save(session.State)
			wg.Done()
		}(states[i].TelegramID)
	}

	logrus.WithField("sent_notifications", sentNotifications).Info("sent birthday reminder")

	wg.Wait()
}
