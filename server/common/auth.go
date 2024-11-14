package common

import (
	"fmt"
	"time"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

// jwtCustomClaims are custom claims extending default ones.
// See https://github.com/golang-jwt/jwt for more examples
type jwtCustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func LoginGetUserToken(user entity.User, key string, ttl int) (string, error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),                                       // 签发时间
			ExpiresAt: time.Now().Add(time.Second * time.Duration(ttl)).Unix(), // 过期时间
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return t, nil
}

var ErrTokenNotProvided = fmt.Errorf("token not provided")
var ErrTokenUserNotFound = fmt.Errorf("user not found")
var ErrTokenInvalidFromDate = fmt.Errorf("token is invalid starting from a certain date")

func GetUserByReq(app *core.App, c *fiber.Ctx) (entity.User, error) {
	uc, err := DecryptAuth(app.Conf().TokenSecret, c.Get(fiber.HeaderAuthorization))
	if err != nil {
		log.Errorf("GetUserByReq|DecryptAuth token:%s err:%v", c.Get(fiber.HeaderAuthorization), err)
		return entity.User{}, err
	}

	email := fmt.Sprintf("%s@9466.com", uc.UserID)
	users := app.Dao().FindUsersByEmail(email)
	if len(users) == 0 {
		return entity.User{}, ErrTokenNotProvided
	}

	return users[0], nil
}
