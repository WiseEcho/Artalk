package common

import (
	"errors"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/i18n"
	"github.com/gofiber/fiber/v2"
)

func AdminGuard(app *core.App, handler fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, isAdmin := CheckIsAdminReq(app, c)
		if !isAdmin {
			return RespError(c, 403, i18n.T("Admin access required"), Map{"need_login": true})
		}

		return handler(c)
	}
}

func LoginGuard(app *core.App, handler func(*fiber.Ctx, entity.User) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := GetUserByReq(app, c)
		if err != nil {
			msg := i18n.T("Login required")
			if errors.Is(err, ErrTokenInvalidFromDate) {
				msg = i18n.T("Your authentication token has expired. Please try signing in again.")
			}
			return RespError(c, 401, msg, Map{"need_auth_login": true})
		}
		return handler(c, user)
	}
}

func CheckIsAdminReq(app *core.App, c *fiber.Ctx) (uint, bool) {
	user, err := GetUserByReq(app, c)
	if err != nil {
		return 0, false
	}

	return user.ID, !user.IsEmpty() && user.IsAdmin
}
