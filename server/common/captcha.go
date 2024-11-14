package common

import (
	"fmt"

	"github.com/artalkjs/artalk/v2/internal/captcha"
	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/server/middleware/limiter"
	"github.com/gofiber/fiber/v2"
)

func GetLimiter(c *fiber.Ctx) (lmt *limiter.Limiter, err error) {
	l, ok := c.Locals("limiter").(*limiter.Limiter)
	if l == nil || !ok {
		return nil, RespError(c, 500, "limiter is not initialize, but middleware is used")
	}
	return l, nil
}

func NewCaptchaChecker(app *core.App, c *fiber.Ctx) captcha.Checker {
	user, _ := GetUserByReq(app, c)
	return captcha.NewCaptchaChecker(&captcha.CheckerConf{
		CaptchaConf: app.Conf().Captcha,
		User: captcha.User{
			ID: fmt.Sprint(user.ID),
			IP: c.IP(),
		},
	})
}

func LimiterGuard(app *core.App, handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return handler(c)
	}
}
