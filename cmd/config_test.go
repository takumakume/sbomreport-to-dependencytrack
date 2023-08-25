package cmd

import (
	"reflect"
	"testing"
)

func Test_newConfig(t *testing.T) {
	type args struct {
		baseURL        string
		apiKey         string
		projectName    string
		projectVersion string
		projectTags    []string
		timeout        int
	}
	tests := []struct {
		name    string
		args    args
		want    *config
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				baseURL:        "http://localhost:8080",
				apiKey:         "apiKey",
				projectName:    "projectName",
				projectVersion: "projectVersion",
				projectTags:    []string{"tag1", "tag2"},
				timeout:        10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newConfig(tt.args.baseURL, tt.args.apiKey, tt.args.projectName, tt.args.projectVersion, tt.args.projectTags, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("newConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
