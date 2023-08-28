package core

import (
	"github.com/ArtalkJS/Artalk/internal/email"
	"github.com/ArtalkJS/Artalk/internal/entity"
	"github.com/ArtalkJS/Artalk/internal/notify_pusher"
)

var _ Service = (*NotifyService)(nil)

type NotifyService struct {
	app    *App
	pusher *notify_pusher.NotifyPusher
}

func NewNotifyService(app *App) *NotifyService {
	return &NotifyService{app: app}
}

func (s *NotifyService) Init() error {
	s.pusher = notify_pusher.NewNotifyPusher(&notify_pusher.NotifyPusherConf{
		AdminNotifyConf: s.app.Conf().AdminNotify,
		Dao:             s.app.Dao(),
		EmailPush: func(notify *entity.Notify) error {
			AppService[*EmailService](s.app).AsyncSend(notify)
			return nil
		},
		EmailRender: func() *email.Render {
			return AppService[*EmailService](s.app).GetRender()
		},
	})

	return nil
}

func (s *NotifyService) Dispose() error {
	s.pusher = nil

	return nil
}

func (s *NotifyService) Push(comment *entity.Comment, pComment *entity.Comment) error {
	s.pusher.Push(comment, pComment)
	return nil
}
