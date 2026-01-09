package medication

import (
	"errors"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
)

const (
	maxCommentaryLength = 1000
)

var (
	ErrInvalidCommentary = errors.New("invalid commentary")
)

// Commentary is a Value Object representing additional notes or comments about a medication.
type Commentary string

// NewCommentary creates validated medication commentary.
func NewCommentary(commentary string) (Commentary, error) {
	err := errors.Join(
		validation.MaxLength(commentary, maxCommentaryLength),
	)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrInvalidCommentary, err)
	}
	return Commentary(commentary), nil
}

// String implements Stringer interface for Commentary.
func (c Commentary) String() string {
	return string(c)
}
