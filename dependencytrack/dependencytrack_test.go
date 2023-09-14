package dependencytrack

import (
	"errors"
	"testing"

	dtrack "github.com/DependencyTrack/client-go"
)

func TestIsNotFound(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "project not found",
			args: args{
				err: ErrProjectNotFound,
			},
			want: true,
		},
		{
			name: "api error 404",
			args: args{
				err: &dtrack.APIError{
					StatusCode: 404,
				},
			},
			want: true,
		},
		{
			name: "api error not 404",
			args: args{
				err: &dtrack.APIError{
					StatusCode: 500,
				},
			},
			want: false,
		},
		{
			name: "other error",
			args: args{
				err: errors.New("error"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.args.err); got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}
