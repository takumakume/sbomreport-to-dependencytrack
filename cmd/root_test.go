package cmd

import (
	"testing"
)

func TestGetBOM(t *testing.T) {
	tests := []struct {
		name           string
		sbomReportJSON []byte
		want           []byte
		wantErr        bool
	}{
		{
			name: "valid sbomReportJSON",
			sbomReportJSON: []byte(`{
                "kind": "SbomReport",
                "apiVersion": "aquasecurity.github.io/v1alpha1",
                "metadata": {
					"name": "test"
				},
                "report": {
                    "bom": {
                        "components": []
                    }
                }
            }`),
			want:    []byte(`{"components":[]}`),
			wantErr: false,
		},
		{
			name:           "invalid sbomReportJSON",
			sbomReportJSON: []byte(`invalid`),
			want:           nil,
			wantErr:        true,
		},
		{
			name: "invalid kind",
			sbomReportJSON: []byte(`{
                "kind": "InvalidKind",
                "apiVersion": "aquasecurity.github.io/v1alpha1",
                "metadata": {
					"name": "test"
				},
                "report": {
                    "bom": {
                        "components": []
                    }
                }
            }`),
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid report",
			sbomReportJSON: []byte(`{
                "kind": "SbomReport",
                "apiVersion": "aquasecurity.github.io/v1alpha1",
                "metadata": {
					"name": "test"
				}
            }`),
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid bom",
			sbomReportJSON: []byte(`{
                "kind": "SbomReport",
                "apiVersion": "aquasecurity.github.io/v1alpha1",
				"metadata": {
					"name": "test"
				},
                "report": {}
            }`),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBOM(tt.sbomReportJSON)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBOM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != string(tt.want) {
				t.Errorf("getBOM() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
