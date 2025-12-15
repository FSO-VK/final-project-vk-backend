package main

import (
	"context"
	"errors"
	httpErr "net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/application"
	generator "github.com/FSO-VK/final-project-vk-backend/internal/planning/application/generate_record"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/application/notify"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/config"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/daemon"
	medClient "github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/medication_client"
	notifyClient "github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/notification_client"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/storage/memory"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/presentation/http"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/configuration"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	auth "github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
	"github.com/sirupsen/logrus"
)

const (
	creationShift  = 24 * time.Hour
	batchSize      = 1000
	tickerInterval = 24 * time.Hour
	// in docker container need to use UTC time.
	timeStart             = 21*time.Hour + 0*time.Minute
	notificationsInterval = 1 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	l.SetReportCaller(true)
	l.SetLevel(logrus.DebugLevel)
	logger := logrus.NewEntry(l)

	confPath, err := configuration.ReadConfigPathFlag("config/planning-conf.yaml")
	if err != nil {
		logger.Fatal(err)
	}
	var conf config.Config
	if err := configuration.KoanfLoad(confPath, &conf); err != nil {
		logger.Fatal(err)
	}

	planRepo := memory.NewPlanStorage()
	recordsRepo := memory.NewRecordStorage()
	medicationClient := medClient.NewMedicationClient(conf.Medication, logger)

	// Service and daemon for generating records
	generateRecordsService := generator.NewGenerateRecordService(recordsRepo, planRepo)
	daemonRecordsGenerator := daemon.NewDaemon(tickerInterval, timeStart, logger)

	// Service and daemon for notifications
	notificationProvider := notifyClient.NewNotificationClient(conf.Notification, logger)
	notifyService := notify.NewIntakeNotificationService(
		recordsRepo,
		planRepo,
		*notificationProvider,
	)
	daemonNotify := daemon.NewDaemon(notificationsInterval, timeStart, logger)

	// Initial generation
	if err := generateRecordsService.GenerateRecordsForDay(ctx, batchSize, creationShift); err != nil {
		logger.Fatal(err)
	}

	validator := validator.NewValidationProvider()
	app := &application.PlanningApplication{
		GetAllPlans: application.NewGetAllPlansService(planRepo, validator),
		GetPlan:     application.NewGetPlanService(planRepo, validator),
		AddPlan: application.NewAddPlanService(
			planRepo,
			generateRecordsService,
			validator,
			medicationClient,
			creationShift,
		),
		ShowSchedule: application.NewShowScheduleService(
			planRepo,
			recordsRepo,
			validator,
			medicationClient,
			creationShift,
		),
		DeletePlan: application.NewFinishPlanService(planRepo, validator),
	}
	planningHandlers := http.NewHandlers(app, logger)

	authChecker := auth.NewHTTPAuthChecker(conf.Auth, logger)
	authMw := httputil.NewAuthMiddleware(authChecker)

	router := http.Router(planningHandlers, authMw)
	server := http.NewGINServer(&conf.Server, logger)
	server.Router(router)

	var wg sync.WaitGroup

	// Shutdown goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		logger.Info("Shutdown signal received")
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Errorf("graceful shutdown failed: %v", err)
		}
	}()

	// Daemon goroutine - generate records
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Daemon started (notifications generation)")
		daemonNotify.Run(ctx, func(ctx context.Context) error {
			return generateRecordsService.GenerateRecordsForDay(ctx, batchSize, creationShift)
		})
	}()

	// Daemon goroutine - send notifications
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Daemon started (records generation)")
		daemonRecordsGenerator.Run(ctx, func(ctx context.Context) error {
			return notifyService.GenerateNotifications(ctx)
		})
	}()

	// Start server
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, httpErr.ErrServerClosed) {
		logger.Fatal(err)
	}

	wg.Wait()
	logger.Info("Server stopped")
}
