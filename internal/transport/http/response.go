package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/FSO-VK/final-project-vk-backend/pkg/api"
)

var ErrResponseIsNil = errors.New("response is nil")

func WriteJSON[T any](w http.ResponseWriter, response *api.Response[T]) error {
	w.Header().Set("Content-Type", "application/json")

	if response == nil {
		return ErrResponseIsNil
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}
