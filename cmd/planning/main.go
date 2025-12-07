package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	generator "github.com/FSO-VK/final-project-vk-backend/internal/planning/application/generate_record"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/config"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/daemon"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/storage/memory"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/configuration"
	"github.com/sirupsen/logrus"
)

const (
	creation_shift  = 24 * time.Hour
	batch_size      = 1000
	ticker_interval = 24 * time.Hour
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
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

	confPath, err := configuration.ReadConfigPathFlag("config/planning-conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	var conf config.Config
	err = configuration.KoanfLoad(confPath, &conf)
	if err != nil {
		logger.Fatal(err)
	}
	planRepo := memory.NewPlanStorage()
	recordsRepo := memory.NewRecordStorage()

	generateRecordsService := generator.NewGenerateRecordService(
		recordsRepo,
		planRepo,
	)
	daemonRecordsGenerator := daemon.NewDaemon(
		ticker_interval,
	)

	if err := generateRecordsService.GenerateRecordsForDay(ctx, batch_size, creation_shift); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Daemon started")
	go daemonRecordsGenerator.Run(ctx, func(ctx context.Context) error {
		return generateRecordsService.GenerateRecordsForDay(ctx, batch_size, creation_shift)
	})
	<-ctx.Done()
	logger.Info("Server stopped")
}
