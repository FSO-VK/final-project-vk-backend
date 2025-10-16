// Package config is a package for configuration for the medication service.
package config

import "github.com/FSO-VK/final-project-vk-backend/internal/medication/transport/http"

// Config is a configuration for the medication service.
type Config struct {
	Server http.ServerConfig
}
