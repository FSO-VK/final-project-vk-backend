// Package config is a package for configuration for the planning service.
package config

import (
	generator "github.com/FSO-VK/final-project-vk-backend/internal/planning/application/generate_record"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/infrastructure/daemon"
	"github.com/FSO-VK/final-project-vk-backend/internal/planning/presentation/http"
	auth "github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
)

// Config is a configuration for the planning service.
type Config struct {
	Server         http.ServerConfig
	Auth           auth.ClientConfig
	GenerateRecord generator.ClientConfig
	GenerateDaemon daemon.ClientConfig
}
