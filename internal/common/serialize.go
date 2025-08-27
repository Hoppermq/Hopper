package common

import (
	"bytes"
	"encoding/gob"

	"github.com/hoppermq/hopper/pkg/domain"
)

// Constraint define item that could be serializable or not.
type Constraint interface {
	comparable
	domain.Serializable
}

// Serializable is a generic type for Serializable items used in the constraints.
type Serializable[T Constraint] struct {
}

// Serialize transform a value into a []byte.
func Serialize[T any](d T) ([]byte, error) {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(d)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// Deserialize transform a []byte into a value.
func Deserialize[T any](d []byte) (T, error) {
	var res T

	buff := bytes.NewReader(d)
	decoder := gob.NewDecoder(buff)
	err := decoder.Decode(&res)

	return res, err
}
