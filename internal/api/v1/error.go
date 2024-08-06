package v1

import (
	"errors"
	"github.com/gofiber/fiber/v3"
)

var (
	ErrInvalidAuthHeader = errors.New("invalid authorization header")
	ErrInvalidAuthToken  = errors.New("invalid authorization token")
)

func errorResponse(c fiber.Ctx, status int, err error) {
	var HTTPError *fiber.Error
	if ok := errors.As(err, &HTTPError); !ok {
		err = fiber.NewError(status, err.Error())
	}
	c.Response().SetStatusCode(status)
	_ = c.JSON(err)
}
