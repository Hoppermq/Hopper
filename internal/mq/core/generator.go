package core

import "github.com/google/uuid"

// GenerateIdentifier generates a new unique identifier using UUID
func GenerateIdentifier() string {
	return uuid.NewString()
}
