// Package config is a package for configuration for the planning service.
package config

import (
	medication "github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/medication_client"
	notification "github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/notification_client"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/presentation/http"
	auth "github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
)

// Config is a configuration for the planning service.
type Config struct {
	Server       http.ServerConfig
	Auth         auth.ClientConfig
	Medication   medication.ClientConfig
	Notification notification.ClientConfig
}
