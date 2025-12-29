package main

import (
	"context"
	"errors"
	httpErr "net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/FSO-VK/final-project-vk-backend/internal/planning/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/config"
	medClient "github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/medication_client"
	notifyProvider "github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/notification"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/storage/memory"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/presentation/http"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/configuration"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/daemon"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	notifyClient "github.com/FSO-VK/final-project-vk-backend/internal/utils/notification_client"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	auth "github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
	"github.com/sirupsen/logrus"
)

const (
	creationShift         = 24 * time.Hour
	batchSize             = 1000
	tickerInterval        = 24 * time.Hour
	notificationsInterval = 1 * time.Minute
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

	now := time.Now()

	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		logger.Fatal(err)
	}

	nowUTC := time.Now().In(loc)

	midnight := time.Date(
		nowUTC.Year(), nowUTC.Month(), nowUTC.Day(),
		0, 0, 0, 0,
		loc,
	).Add(24 * time.Hour)

	quickStart := now.Add(2 * time.Minute)

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
	generateRecordsService := application.NewGenerateRecordService(recordsRepo, planRepo)
	daemonRecordsGenerator := daemon.NewDaemon(tickerInterval, midnight, logger)

	// Service and daemon for intake notifications
	notificationProvider := notifyClient.NewNotificationClient(conf.Notification, logger)
	notificationAdapter := notifyProvider.NewNotificationProvider(notificationProvider)
	intakeNotificationService := application.NewIntakeNotificationService(
		recordsRepo,
		planRepo,
		notificationAdapter,
		medicationClient,
	)

	daemonIntakeNotification := daemon.NewDaemon(notificationsInterval, quickStart, logger)

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
		logger.Info("Daemon started (records generation)")
		daemonRecordsGenerator.Run(ctx, func(ctx context.Context) error {
			return generateRecordsService.GenerateRecordsForDay(ctx, batchSize, creationShift)
		})
	}()

	// Daemon goroutine - send intake notifications
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Daemon started (intake notifications generation)")
		daemonIntakeNotification.Run(ctx, func(ctx context.Context) error {
			return intakeNotificationService.GenerateIntakeNotifications(ctx)
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
