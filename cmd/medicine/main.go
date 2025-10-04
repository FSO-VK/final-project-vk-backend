package main

import (
	"github.com/FSO-VK/final-project-vk-backend/internal/medicine/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/medicine/storage/memory"
	"github.com/FSO-VK/final-project-vk-backend/internal/medicine/transport/http"
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

	medicineRepo := memory.NewMedicineStorage()

	medicineService := application.NewMedicineServiceProvider(medicineRepo)

	medicineHandlers := http.NewHandlers(medicineService, logger)
	router := http.Router(medicineHandlers)

	server := http.NewHTTPServer(":8080", logger)
	server.Router(router)

	_ = server.ListenAndServe()
}
