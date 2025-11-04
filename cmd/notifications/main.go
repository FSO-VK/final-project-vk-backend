package main

import (
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/infrastructure/config"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/infrastructure/datamatrix"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/infrastructure/storage/memory"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/presentation/http"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/configuration"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	auth "github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
	webpush "github.com/SherClockHolmes/webpush-go"
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

	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		// TODO: Handle error
	}
	fmt.Println("!", privateKey, publicKey)

	notificationsRepo := memory.NewNotificationsStorage()
	validator := validator.NewValidationProvider()
	dataMatrixClient := datamatrix.NewDataMatrixAPI(
		conf.Scan,
		logger,
	)
	dataMatrixCache := memory.NewDataMatrixStorage()
	notificationsBoxRepo := memory.NewNotificationsBoxStorage()
	app := &application.NotificationsApplication{
		// GetMedicationBox: application.NewGetMedicationBoxService(
		// 	medicationRepo, medicationBoxRepo, validator),
		// AddMedication: application.NewAddMedicationService(
		// 	medicationRepo, medicationBoxRepo, validator),
		// UpdateMedication: application.NewUpdateMedicationService(
		// 	medicationRepo, medicationBoxRepo, validator),
		// DeleteMedication: application.NewDeleteMedicationService(
		// 	medicationRepo, medicationBoxRepo, validator),
		// DataMatrixInformation: application.NewDataMatrixInformationService(
		// 	dataMatrixClient,
		// 	dataMatrixCache,
		// 	validator,
		// ),
	}
	notificationsHandlers := http.NewHandlers(app, logger)

	authChecker := auth.NewHTTPAuthChecker(conf.Auth, logger)

	authMw := httputil.NewAuthMiddleware(authChecker)

	router := http.Router(notificationsHandlers, authMw)

	server := http.NewHTTPServer(&conf.Server, logger)
	server.Router(router)

	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}
