package v1

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"news_api/internal/service"
	"os"
	"strings"
)

const bearerPrefix = "Bearer "

type authMiddleware struct {
	auth service.Auth
}

func (h *authMiddleware) authHandler(c fiber.Ctx) error {
	token, ok := parseToken(c.Request())
	if !ok {
		errorResponse(c, http.StatusUnauthorized, ErrInvalidAuthHeader)
		return nil
	}
	if !h.auth.ValidateToken(token) {
		errorResponse(c, http.StatusForbidden, ErrInvalidAuthToken)
		return nil
	}
	return c.Next()
}

func parseToken(r *fasthttp.Request) (string, bool) {
	header := string(r.Header.Peek(fiber.HeaderAuthorization))
	if header == "" {
		return "", false
	}
	token := strings.Split(header, bearerPrefix)
	if len(token) != 2 {
		return "", false
	}
	return token[1], true
}

func LoggingMiddleware(h *fiber.App, output string) {
	config := logger.Config{
		Format: `{"time":"${time_rfc3339}", "method":"${method}","uri":"${uri}", "status":${status}, "error":"${error}"}` + "\n",
	}
	if output == "stdout" {
		config.Output = os.Stdout
	} else {
		file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.Fatal(err)
		}
		config.Output = file
	}
	h.Use(logger.New(config))
}
