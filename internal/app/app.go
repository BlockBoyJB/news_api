package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	v1 "news_api/internal/api/v1"
	"news_api/internal/repo"
	"news_api/internal/service"
	"news_api/pkg/validator"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"news_api/config"
	"news_api/pkg/httpserver"
	"news_api/pkg/postgres"
	"os"
)

func Run() {
	// config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	setLogger(cfg.Log.Level, cfg.Log.Output)

	// postgres
	pg, err := postgres.NewPG(cfg.PG.Url, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	d := &service.ServicesDependencies{
		Repos:      repo.NewRepositories(pg),
		PrivateKey: cfg.JWT.PrivateKey,
		PublicKey:  cfg.JWT.PublicKey,
	}
	services := service.NewServices(d)

	// validator for incoming messages
	v, err := validator.NewValidator()
	if err != nil {
		log.Fatalf("Initializing handler validator error: %s", err)
	}

	// handler for incoming messages
	h := fiber.New()
	v1.LoggingMiddleware(h, cfg.Log.Output)
	v1.NewRouter(h, services, v)

	httpServer := httpserver.NewServer(h.Handler(), httpserver.Port(cfg.HTTP.Port))

	log.Infof("App started! Listening port %s", cfg.HTTP.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app run, signal " + s.String())

	case err = <-httpServer.Notify():
		log.Errorf("/app/run http server notify error: %s", err)
	}
	// graceful shutdown
	if err = httpServer.Shutdown(); err != nil {
		log.Errorf("/app/run http server shutdown error: %s", err)
	}

	log.Infof("App shutdown with exit code 0")

}

// loading environment params from .env
func init() {
	if _, ok := os.LookupEnv("HTTP_PORT"); !ok {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("load env file error: %s", err)
		}
	}
}
