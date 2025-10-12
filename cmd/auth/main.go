package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/FSO-VK/final-project-vk-backend/internal/auth/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/auth/infrastructure/config"
	"github.com/FSO-VK/final-project-vk-backend/internal/auth/infrastructure/storage/memory"
	"github.com/FSO-VK/final-project-vk-backend/internal/auth/presentation/http"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/configuration"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/password"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	"github.com/sirupsen/logrus"
)

func main() {
	l := logrus.New()
	l.SetFormatter(
		&logrus.TextFormatter{
			ForceColors:            true,
			FullTimestamp:          true,
			DisableLevelTruncation: true,
			PadLevelText:           true,
		},
	)
	l.SetReportCaller(true)
	l.SetLevel(logrus.DebugLevel)
	logger := logrus.NewEntry(l)

	var conf config.Config
	err := configuration.KoanfLoad("config/auth-conf.yaml", &conf)
	if err != nil {
		logger.Fatal(err)
	}

	validator := validator.NewValidationProvider()
	credentialRepo := memory.NewCredentialStorage()
	sessionRepo := memory.NewSessionStorage()
	hasher := password.NewPasswordHasherProvider()

	app := &application.AuthApplication{
		LoginByEmail: application.NewLoginByEmailService(
			credentialRepo,
			sessionRepo,
			validator,
			hasher,
		),
		Logout: application.NewLogoutService(
			sessionRepo,
			validator,
		),
		CheckAuth: application.NewCheckAuthService(
			sessionRepo,
			validator,
		),
		Registration: application.NewRegistrationService(
			credentialRepo,
			sessionRepo,
			validator,
			hasher,
		),
	}

	handlers := http.NewAuthHandlers(
		app,
		logger,
	)

	router := http.NewRouter(handlers)

	server := http.NewServerHTTP(conf.Server, router, logger)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

		<-sigint

		logger.Info("shutting down server...")

		err := server.Shutdown()
		if err != nil {
			logger.Errorf("server shutdown: %v", err)
		}
	}()

	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}

	wg.Wait()

	logger.Info("server stopped")
}
