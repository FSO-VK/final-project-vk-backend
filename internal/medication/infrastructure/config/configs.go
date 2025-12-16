// Package config is a package for configuration for the medication service.
package config

import (
	dataMatrixClient "github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/datamatrix"
	llm "github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/llm_chat_bot"
	vidalclient "github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/vidal/client"
	vidalstorage "github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/vidal/storage/mongo"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/presentation/http"
	notification "github.com/FSO-VK/final-project-vk-backend/internal/utils/notification_client"
	auth "github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
	"github.com/FSO-VK/final-project-vk-backend/pkg/llm/gigachat"
)

// Config is a configuration for the medication service.
type Config struct {
	Server       http.ServerConfig
	Internal     http.ServerConfig
	Scan         dataMatrixClient.ClientConfig
	Auth         auth.ClientConfig
	Vidal        vidal
	Gigachat     gigachat.ClientConfig
	Assistant    llm.InstructionAssistantConfig
	Notification notification.ClientConfig
}

type vidal struct {
	Client  vidalclient.Config
	Storage vidalstorage.Config
}
