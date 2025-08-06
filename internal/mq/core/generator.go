package core

import "github.com/google/uuid"

func GenerateIdentifier() string {
	return uuid.NewString()
}
