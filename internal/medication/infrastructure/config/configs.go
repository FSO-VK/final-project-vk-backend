// Package config is a package for configuration for the medication service.
package config

import (
	dataMatrixClient "github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/datamatrix"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/presentation/http"
	auth "github.com/FSO-VK/final-project-vk-backend/pkg/auth/client"
)

// Config is a configuration for the medication service.
type Config struct {
	Server http.ServerConfig
	Scan   dataMatrixClient.ClientConfig
	Auth   auth.ClientConfig
}
