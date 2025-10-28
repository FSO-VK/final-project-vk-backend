// Package config is a package for configuration for the medication service.
package config

import (
	dataMatrixClient "github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/client_data_matrix"
	"github.com/FSO-VK/final-project-vk-backend/internal/medication/presentation/http"
)

// Config is a configuration for the medication service.
type Config struct {
	Server http.ServerConfig
	Scan   dataMatrixClient.ClientConfig
}
