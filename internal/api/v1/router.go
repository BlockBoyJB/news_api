package v1

import (
	"github.com/gofiber/fiber/v3"
	"net/http"
	"news_api/internal/service"
	"news_api/pkg/validator"
)

func NewRouter(h *fiber.App, services *service.Services, validator validator.Validator) {
	h.Get("/ping", ping)

	auth := authMiddleware{auth: services.Auth}
	h.Get("/token", auth.getToken)

	v1 := h.Group("/api/v1", auth.authHandler)
	newNewsRouter(v1.Group("/news"), services.News, validator)
}

func ping(c fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

// Вообще в микросервисе этого не должно быть, но так как я разрабатываю его изолированно, то можно
func (h *authMiddleware) getToken(c fiber.Ctx) error {

	type response struct {
		Token string `json:"token"`
	}
	token, err := h.auth.CreateToken()
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, err)
		return err
	}
	return c.JSON(response{Token: token})
}
