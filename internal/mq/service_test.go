package mq

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateService(t *testing.T) {
	t.Parallel()
	type args struct {
		opts []Option
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ValueAssertionFunc
	}{
		{
			name: "Create_Service",
			args: args{
				opts: []Option{
					WithTCP(),
					WithLogger(
						slog.New(
							slog.NewJSONHandler(os.Stdout, nil),
						),
					),
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := New(tt.args.opts...)
			assert.Equal(t, tt.want, got != nil)
		})
	}

}
