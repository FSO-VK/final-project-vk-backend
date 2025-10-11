package configuration

import (
	"fmt"
	"os"

	"github.com/drone/envsubst"
	"github.com/joho/godotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
)

func KoanfLoad(path string, configStruct interface{}) error {
	k := koanf.New(".")

	// It is possible to ignore error caused by .env because it's optional.
	_ = godotenv.Load()

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read yaml file: %w", err)
	}

	processedFile, err := envsubst.EvalEnv(string(fileBytes))
	if err != nil {
		return fmt.Errorf("failed to process env substitution: %w", err)
	}

	err = k.Load(rawbytes.Provider([]byte(processedFile)), yaml.Parser())
	if err != nil {
		return fmt.Errorf("failed to load yaml: %w", err)
	}

	err = k.Unmarshal("", configStruct)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return nil
}
