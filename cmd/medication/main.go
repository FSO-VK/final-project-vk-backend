package main

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/config"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/storage/memory"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/transport/http"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/configuration"
	"github.com/sirupsen/logrus"
)

func main() {
	// TODO: configure logger via config file
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

	medicineRepo := memory.NewMedicineStorage()

	medicineService := application.NewMedicineServiceProvider(medicineRepo)

	medicineHandlers := http.NewHandlers(medicineService, logger)
	router := http.Router(medicineHandlers)

	server := http.NewHTTPServer(&conf.Server, logger)
	server.Router(router)

	_ = server.ListenAndServe()
}
