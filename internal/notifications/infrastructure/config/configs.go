// Package config is a package for configuration for the notifications service.
package config

import (
	client "github.com/FSO-VK/final-project-vk-backend/internal/notifications/infrastructure/push_provider"
	"github.com/FSO-VK/final-project-vk-backend/internal/notifications/presentation/http"
	auth "github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
)

// Config is a configuration for the notifications service.
type Config struct {
	Server     http.ServerConfig
	PushClient client.PushClient
	Auth       auth.ClientConfig
}
