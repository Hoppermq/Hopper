package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testNewConfig(t *testing.T) {
	type args struct {
		env string
	}

	tests := []struct {
		name    string
		arg     args
		want    bool
		wantErr assert.ValueAssertionFunc
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("HOPPER_ENV", test.arg.env)
			got, err := New("hopper")
			test.wantErr(t, err)
			assert.Equal(t, test.want, got != nil)
		})
	}
}
