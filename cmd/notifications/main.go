package main

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/infrastructure/client"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/infrastructure/config"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/infrastructure/storage/memory"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/presentation/http"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/configuration"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	auth "github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
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

	confPath, err := configuration.ReadConfigPathFlag("config/notifications-conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	var conf config.Config
	err = configuration.KoanfLoad(confPath, &conf)
	if err != nil {
		logger.Fatal(err)
	}

	subscriptionsRepo := memory.NewSubscriptionsStorage()
	validator := validator.NewValidationProvider()

	pushProvider := client.NewPushNotificationProvider(
		conf.PushClient,
		logger,
	)

	app := &application.NotificationsApplication{
		GetVapidPublicKey: application.NewGetVapidPublicKeyService(
			application.PublicKey(conf.PushClient.VapidPublicKey), validator),
		CreateSubscription: application.NewCreateSubscriptionService(
			subscriptionsRepo,
			validator,
		),
		DeleteSubscription: application.NewDeleteSubscriptionService(
			subscriptionsRepo,
			validator,
		),
		SendNotification: application.NewSendNotificationService(
			subscriptionsRepo,
			pushProvider,
			validator,
		),
	}
	notificationsHandlers := http.NewHandlers(app, logger)

	authChecker := auth.NewHTTPAuthChecker(conf.Auth, logger)

	authMw := httputil.NewAuthMiddleware(authChecker)

	router := http.Router(notificationsHandlers, authMw)

	server := http.NewGINServer(&conf.Server, logger)
	server.Router(router)
	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}
