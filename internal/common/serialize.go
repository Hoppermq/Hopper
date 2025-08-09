package common

import (
	"bytes"
	"encoding/gob"

	"github.com/hoppermq/hopper/pkg/domain"
)

type Constraint interface {
	comparable
	domain.Serializable
}

type Serializable[T Constraint] struct {
}

func Serialize[T any](d T) ([]byte, error) {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(d)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func Deserialize[T any](d []byte) (T, error) {
	var res T

	buff := bytes.NewReader(d)
	decoder := gob.NewDecoder(buff)
	err := decoder.Decode(&res)

	return res, err
}
