package medication

import (
	"errors"
	"fmt"

	"github.com/FSO-VK/final-project-vk-backend/pkg/validation"
)

const maxGroupNameLength = 100

var ErrInvalidGroup = errors.New("invalid group")

// Group is a Value Object representing
// the therapeutic or pharmacological group of a medication.
type Group string

// NewGroup creates validated medication group.
func NewGroup(group string) (Group, error) {
	err := validation.MaxLength(group, maxGroupNameLength)
	if err != nil {
		return Group(""), fmt.Errorf("%w: %w", ErrInvalidGroup, err)
	}
	return Group(group), nil
}

// String implements stringer interface for Group.
func (g Group) String() string {
	return string(g)
}
