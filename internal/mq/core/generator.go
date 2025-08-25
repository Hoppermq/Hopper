package core

import (
	"github.com/google/uuid"

	"github.com/hoppermq/hopper/pkg/domain"
)

// GenerateIdentifier generates a new unique identifier using UUID
func GenerateIdentifier() domain.ID {
	return domain.ID(uuid.NewString())
}
