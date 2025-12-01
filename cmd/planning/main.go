package main

import (
	"context"
	"os"
	"os/signal"

	generator "github.com/FSO-VK/final-project-vk-backend/internal/planning/application/generate_record"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/config"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/storage/memory"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/configuration"
	"github.com/sirupsen/logrus"
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
	ticker := generator.NewTicker(conf.GenerateRecord.TickerInterval)

	generateRecordsService := generator.NewGenerateRecordService(
		conf.GenerateRecord,
		recordsRepo,
		planRepo,
		ticker,
	)
	if err := generateRecordsService.GenerateRecordsForDay(ctx); err != nil {
		logger.Fatal(err)
	}
	logger.Info("Daemon started")
	go generateRecordsService.Run(ctx)
	<-ctx.Done()
	logger.Info("Server stopped")
}
