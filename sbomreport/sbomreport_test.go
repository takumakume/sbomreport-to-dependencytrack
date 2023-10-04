package sbomreport

import (
	"fmt"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	sbomReportV1alpha1, err := os.ReadFile("../testdata/v1alpha1.json")
	if err != nil {
		t.Fatal(err)
	}
	sbomReportV1alpha1WithVerb, err := os.ReadFile("../testdata/v1alpha1_with_verb.json")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		rawJSON []byte
	}
	tests := []struct {
		name     string
		args     args
		wantVerb string
		wantErr  bool
	}{
		{
			name: "success v1alpha1",
			args: args{
				rawJSON: sbomReportV1alpha1,
			},
			wantVerb: "update",
			wantErr:  false,
		},
		{
			name: "success v1alpha1 with verb",
			args: args{
				rawJSON: sbomReportV1alpha1WithVerb,
			},
			wantVerb: "delete",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.rawJSON)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.bom == nil {
				t.Errorf("New() bom is nil")
			}
			if got.verb != tt.wantVerb {
				t.Errorf("New() verb = %v, wantVerb %v", got.verb, tt.wantVerb)
			}
		})
	}
}

func TestIsErrNotSBOMReport(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "error is SbomReportError",
			args: args{
				err: ErrNotSBOMReport,
			},
			want: true,
		},
		{
			name: "error is not SbomReportError",
			args: args{
				err: fmt.Errorf("some error"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsErrNotSBOMReport(tt.args.err); got != tt.want {
				t.Errorf("IsErrNotSBOMReport() = %v, want %v", got, tt.want)
			}
		})
	}
}
