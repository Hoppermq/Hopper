package mocks_generator

import (
	"github.com/hoppermq/hopper/pkg/domain"
)

type MockGenerator struct {
	Generator domain.Generator
}

func NewMockGenerator(generator domain.Generator) *MockGenerator {
	m := &MockGenerator{
		generator,
	}

	return m
}

func (m *MockGenerator) ON(s string) {
	m.Generator = func() domain.ID {
		return domain.ID(s)
	}
}
