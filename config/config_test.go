package config

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		baseURL                    string
		apiKey                     string
		projectName                string
		projectVersion             string
		projectTags                []string
		parentName                 string
		parentVersion              string
		dtrackClientTimeoutSec     float64
		sbomUploadTimeoutSec       float64
		sbomUploadCheckIntervalSec float64
	}
	tests := []struct {
		name string
		args args
		want *Config
	}{
		{
			name: "success",
			args: args{
				baseURL:                    "https://example.com",
				apiKey:                     "12345",
				projectName:                "test-project",
				projectVersion:             "1.0.0",
				projectTags:                []string{"tag1", "tag2"},
				parentName:                 "TEST",
				parentVersion:              "1.0.0",
				dtrackClientTimeoutSec:     10,
				sbomUploadTimeoutSec:       30,
				sbomUploadCheckIntervalSec: 1,
			},
			want: &Config{
				BaseURL:                 "https://example.com",
				APIKey:                  "12345",
				ProjectName:             "test-project",
				ProjectVersion:          "1.0.0",
				ProjectTags:             []string{"tag1", "tag2"},
				ParentName:              "TEST",
				ParentVersion:           "1.0.0",
				DtrackClientTimeout:     time.Duration(10) * time.Second,
				SBOMUploadTimeout:       time.Duration(30) * time.Second,
				SBOMUploadCheckInterval: time.Duration(1) * time.Second,
			},
		},
		{
			name: "success tag separator is comma",
			args: args{
				baseURL:                    "https://example.com",
				apiKey:                     "12345",
				projectName:                "test-project",
				projectVersion:             "1.0.0",
				projectTags:                []string{"tag1,tag2"},
				parentName:                 "TEST",
				parentVersion:              "1.0.0",
				dtrackClientTimeoutSec:     10,
				sbomUploadTimeoutSec:       30,
				sbomUploadCheckIntervalSec: 1,
			},
			want: &Config{
				BaseURL:                 "https://example.com",
				APIKey:                  "12345",
				ProjectName:             "test-project",
				ProjectVersion:          "1.0.0",
				ProjectTags:             []string{"tag1", "tag2"},
				ParentName:              "TEST",
				ParentVersion:           "1.0.0",
				DtrackClientTimeout:     time.Duration(10) * time.Second,
				SBOMUploadTimeout:       time.Duration(30) * time.Second,
				SBOMUploadCheckInterval: time.Duration(1) * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.baseURL, tt.args.apiKey, tt.args.projectName, tt.args.projectVersion, tt.args.projectTags, tt.args.parentName, tt.args.parentVersion, tt.args.dtrackClientTimeoutSec, tt.args.sbomUploadTimeoutSec, tt.args.sbomUploadCheckIntervalSec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	type fields struct {
		BaseURL        string
		APIKey         string
		ProjectName    string
		ProjectVersion string
		ProjectTags    []string
		ParentName     string `json:"parentName,omitempty"`
		ParentVersion  string `json:"parentVersion,omitempty"`
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				BaseURL:        "https://example.com",
				APIKey:         "12345",
				ProjectName:    "test-project",
				ProjectVersion: "1.0.0",
				ProjectTags:    []string{"tag1", "tag2"},
				ParentName:     "TEST",
				ParentVersion:  "1.0.0",
			},
			wantErr: false,
		},
		{
			name: "APIKey no set",
			fields: fields{
				BaseURL:        "https://example.com",
				ProjectName:    "test-project",
				ProjectVersion: "1.0.0",
				ProjectTags:    []string{"tag1", "tag2"},
				ParentName:     "TEST",
				ParentVersion:  "1.0.0",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				BaseURL:        tt.fields.BaseURL,
				APIKey:         tt.fields.APIKey,
				ProjectName:    tt.fields.ProjectName,
				ProjectVersion: tt.fields.ProjectVersion,
				ProjectTags:    tt.fields.ProjectTags,
				ParentName:     tt.fields.ParentName,
				ParentVersion:  tt.fields.ParentVersion,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
