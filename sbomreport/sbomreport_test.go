package sbomreport

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	sbomReportV1alpha1, err := os.ReadFile("../testdata/v1alpha1.json")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		rawJSON []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success v1alpha1",
			args: args{
				rawJSON: sbomReportV1alpha1,
			},
			wantErr: false,
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
		})
	}
}
