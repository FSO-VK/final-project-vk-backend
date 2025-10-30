package main

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/config"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/datamatrix"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/storage/memory"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/presentation/http"
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

	confPath, err := configuration.ReadConfigPathFlag("config/medication-conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	var conf config.Config
	err = configuration.KoanfLoad(confPath, &conf)
	if err != nil {
		logger.Fatal(err)
	}
	medicationRepo := memory.NewMedicationStorage()
	validator := validator.NewValidationProvider()
	dataMatrixClient := datamatrix.NewDataMatrixAPI(
		conf.Scan,
		logger,
	)
	dataMatrixCache := memory.NewDataMatrixStorage()
	medicationBoxRepo := memory.NewMedicationBoxStorage()

	app := &application.MedicationApplication{
		GetMedicationBox: application.NewGetMedicationBoxService(
			medicationRepo, medicationBoxRepo, validator),
		AddMedication: application.NewAddMedicationService(
			medicationRepo, medicationBoxRepo, validator),
		UpdateMedication: application.NewUpdateMedicationService(
			medicationRepo, medicationBoxRepo, validator),
		DeleteMedication: application.NewDeleteMedicationService(
			medicationRepo, medicationBoxRepo, validator),
		DataMatrixInformation: application.NewDataMatrixInformationService(
			dataMatrixClient,
			dataMatrixCache,
			validator,
		),
	}

	medicationHandlers := http.NewHandlers(app, logger)

	authChecker := auth.NewHTTPAuthChecker(conf.Auth, logger)

	authMw := httputil.NewAuthMiddleware(authChecker)

	router := http.Router(medicationHandlers, authMw)

	server := http.NewHTTPServer(&conf.Server, logger)
	server.Router(router)

	_ = server.ListenAndServe()
}
