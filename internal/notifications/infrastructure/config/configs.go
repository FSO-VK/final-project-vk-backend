// Package config is a package for configuration for the notifications service.
package config

import (
	client "github.com/FSO-VK/final-project-vk-backend/internal/notifications/infrastructure/client"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/presentation/http"
)

// Config is a configuration for the notifications service.
type Config struct {
	Server http.ServerConfig
	pushClient   client.PushClient
}
