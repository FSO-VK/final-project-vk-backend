package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/FSO-VK/final-project-vk-backend/internal/medication/application"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/config"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/datamatrix"
	instructionAssistant "github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/llm_chat_bot"
	notifyAdapter "github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/notification"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/storage/memory"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/vidal"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/vidal/client"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/vidal/storage/mongo"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/presentation/http"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/configuration"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/daemon"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/httputil"
	notifyClient "github.com/FSO-VK/final-project-vk-backend/internal/utils/notification_client"
	"github.com/FSO-VK/final-project-vk-backend/internal/utils/validator"
	auth "github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
	"github.com/FSO-VK/final-project-vk-backend/pkg/llm/gigachat"
	"github.com/sirupsen/logrus"
)

const (
	notificationsInterval = 24 * time.Hour
	timeDelta             = 7 * 24 * time.Hour
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
	vidalClient := client.NewHTTPClient(conf.Vidal.Client)
	vidalCache, err := mongo.NewStorage(&conf.Vidal.Storage, logger)
	if err != nil {
		logger.Fatal(err)
	}
	medReference := vidal.NewService(
		vidalCache,
		vidalClient,
	)
	medicationBoxRepo := memory.NewMedicationBoxStorage()
	instructionLLMProvider := gigachat.NewGigachatLLMProvider(conf.Gigachat)
	instructionLLM := instructionAssistant.NewLLMChatBot(
		instructionLLMProvider,
		conf.Assistant,
	)

	app := &application.MedicationApplication{
		GetMedicationBox: application.NewGetMedicationBoxService(
			medicationRepo, medicationBoxRepo, validator),
		GetMedicationByID: application.NewGetMedicationByIDService(
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
			medReference,
			validator,
		),
		InstructionAssistant: application.NewInstructionAssistantService(
			medicationRepo,
			instructionLLM,
			medReference,
			validator,
		),
		GetInstructionByMedicationID: application.NewGetInstructionByMedicationIDService(
			medicationRepo,
			medicationBoxRepo,
			medReference,
			validator,
		),
	}

	medicationHandlers := http.NewHandlers(app, logger)

	authChecker := auth.NewHTTPAuthChecker(conf.Auth, logger)

	authMw := httputil.NewAuthMiddleware(authChecker)

	loggingMw := httputil.NewLoggingMiddleware(logger)

	// public router
	router := http.Router(medicationHandlers, authMw, loggingMw)
	server := http.NewHTTPServer(&conf.Server, logger)
	server.Router(router)

	// internal router
	internalRouter := http.InternalRouter(medicationHandlers, loggingMw)
	internalServer := http.NewHTTPServer(&conf.Internal, logger)
	internalServer.Router(internalRouter)

	// daemon expiration notifications
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		logger.Fatal(err)
	}

	nowUTC := time.Now().In(loc)

	noon := time.Date(
		nowUTC.Year(), nowUTC.Month(), nowUTC.Day(),
		12, 0, 0, 0,
		loc,
	).Add(24 * time.Hour)
	daemonExpirationNotification := daemon.NewDaemon(notificationsInterval, noon, logger)
	notificationProvider := notifyClient.NewNotificationClient(conf.Notification, logger)
	notificationAdapter := notifyAdapter.NewAdapter(notificationProvider)
	expirationNotificationService := application.NewExpirationNotificationService(
		medicationRepo,
		medicationBoxRepo,
		notificationAdapter,
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := internalServer.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Daemon started (expiation notifications)")
		daemonExpirationNotification.Run(ctx, func(ctx context.Context) error {
			return expirationNotificationService.GenerateExpirationNotifications(ctx, timeDelta)
		})
	}()

	go func() {
		<-stop
		logger.Info("Servers are shutting down...")
		cancel()

		_ = server.Shutdown(context.Background())
		_ = internalServer.Shutdown(context.Background())
	}()

	wg.Wait()
	logger.Info("All servers stopped")
}
