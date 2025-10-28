package main

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application"
	dataMatrixClient "github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/client_data_matrix"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/config"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/storage/memory"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/presentation/http"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/configuration"
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
	dataMatrixClient := dataMatrixClient.NewDataMatrixAPI(
		conf.Scan,
		logger,
	)
	dataMatrixCache := memory.NewDataMatrixStorage()

	app := &application.MedicationApplication{
		GetMedicationList: application.NewGetMedicationListService(medicationRepo, validator),
		AddMedication:     application.NewAddMedicationService(medicationRepo, validator),
		UpdateMedication:  application.NewUpdateMedicationService(medicationRepo, validator),
		DeleteMedication:  application.NewDeleteMedicationService(medicationRepo, validator),
		DataMatrixInformation: application.NewDataMatrixInformationService(
			dataMatrixClient,
			dataMatrixCache,
			validator,
		),
	}

	medicationHandlers := http.NewHandlers(app, logger)
	router := http.Router(medicationHandlers)

	server := http.NewHTTPServer(&conf.Server, logger)
	server.Router(router)

	_ = server.ListenAndServe()
}
