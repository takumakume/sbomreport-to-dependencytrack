package uploader

import (
	"context"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/takumakume/sbomreport-to-dependencytrack/config"
	"github.com/takumakume/sbomreport-to-dependencytrack/mock"
)

func TestUpload_Run(t *testing.T) {
	sbomReportV1alpha1, err := os.ReadFile("../testdata/v1alpha1.json")
	if err != nil {
		t.Fatal(err)
	}

	sbomReportV1alpha1WithVerb, err := os.ReadFile("../testdata/v1alpha1_with_verb.json")
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDTrack := mock.NewMockDependencyTrackClient(ctrl)

	type mockUploadBOM struct {
		enable         bool
		projectName    string
		projectVersion string
		projectTags    []string
		parentName     string
		parentVersion  string
		err            error
	}
	type mockAddTagsToProject struct {
		enable         bool
		projectName    string
		projectVersion string
		projectTags    []string
		err            error
	}
	type mockDeleteSBOM struct {
		enableDeactivate bool
		enableDelete     bool
		projectName      string
		projectVersion   string
	}
	testCases := []struct {
		name                 string
		config               *config.Config
		input                []byte
		mockUploadBOM        mockUploadBOM
		mockAddTagsToProject mockAddTagsToProject
		mockDeleteSBOM       mockDeleteSBOM
		wantErr              bool
	}{
		{
			name: "success",
			config: &config.Config{
				BaseURL:        "http://localhost:8081",
				APIKey:         "apiKey",
				ProjectName:    "[[.sbomReport.report.artifact.repository]]",
				ProjectVersion: "[[.sbomReport.report.artifact.tag]]",
				ParentName:     "[[.sbomReport.metadata.namespace]]",
				ParentVersion:  "[[.sbomReport.report.artifact.tag]]",
				ProjectTags: []string{
					"test",
					"kube_namespace:[[.sbomReport.metadata.namespace]]",
				},
			},
			input: sbomReportV1alpha1,
			mockUploadBOM: mockUploadBOM{
				enable:         true,
				projectName:    "library/alpine",
				projectVersion: "latest",
				parentName:     "default",
				parentVersion:  "latest",
				projectTags: []string{
					"test",
					"kube_namespace:default",
				},
				err: nil,
			},
			mockAddTagsToProject: mockAddTagsToProject{
				enable:         true,
				projectName:    "library/alpine",
				projectVersion: "latest",
				projectTags: []string{
					"test",
					"kube_namespace:default",
				},
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "delete with action 'ignore' does nothing",
			config: &config.Config{
				BaseURL:          "http://localhost:8081",
				APIKey:           "apiKey",
				ProjectName:      "[[.sbomReport.operatorObject.report.artifact.repository]]",
				ProjectVersion:   "[[.sbomReport.operatorObject.report.artifact.tag]]",
				ParentName:       "[[.sbomReport.operatorObject.metadata.namespace]]",
				ParentVersion:    "[[.sbomReport.operatorObject.report.artifact.tag]]",
				SBOMDeleteAction: "ignore",
				ProjectTags: []string{
					"test",
					"kube_namespace:[[.sbomReport.operatorObject.metadata.namespace]]",
				},
			},
			input:   sbomReportV1alpha1WithVerb,
			wantErr: false,
		},
		{
			name: "delete with action 'deactivate' works",
			config: &config.Config{
				BaseURL:          "http://localhost:8081",
				APIKey:           "apiKey",
				ProjectName:      "[[.sbomReport.operatorObject.report.artifact.repository]]",
				ProjectVersion:   "[[.sbomReport.operatorObject.report.artifact.tag]]",
				ParentName:       "[[.sbomReport.operatorObject.metadata.namespace]]",
				ParentVersion:    "[[.sbomReport.operatorObject.report.artifact.tag]]",
				SBOMDeleteAction: "deactivate",
				ProjectTags: []string{
					"test",
					"kube_namespace:[[.sbomReport.operatorObject.metadata.namespace]]",
				},
			},
			mockDeleteSBOM: mockDeleteSBOM{
				enableDeactivate: true,
				projectName:      "library/alpine",
				projectVersion:   "latest",
			},
			input:   sbomReportV1alpha1WithVerb,
			wantErr: false,
		},
		{
			name: "delete with action 'delete' works",
			config: &config.Config{
				BaseURL:          "http://localhost:8081",
				APIKey:           "apiKey",
				ProjectName:      "[[.sbomReport.operatorObject.report.artifact.repository]]",
				ProjectVersion:   "[[.sbomReport.operatorObject.report.artifact.tag]]",
				ParentName:       "[[.sbomReport.operatorObject.metadata.namespace]]",
				ParentVersion:    "[[.sbomReport.operatorObject.report.artifact.tag]]",
				SBOMDeleteAction: "delete",
				ProjectTags: []string{
					"test",
					"kube_namespace:[[.sbomReport.operatorObject.metadata.namespace]]",
				},
			},
			mockDeleteSBOM: mockDeleteSBOM{
				enableDeactivate: false,
				enableDelete:     true,
				projectName:      "library/alpine",
				projectVersion:   "latest",
			},
			input:   sbomReportV1alpha1WithVerb,
			wantErr: false,
		},
		{
			name: "no tags",
			config: &config.Config{
				BaseURL:        "http://localhost:8081",
				APIKey:         "apiKey",
				ProjectName:    "[[.sbomReport.report.artifact.repository]]",
				ProjectVersion: "[[.sbomReport.report.artifact.tag]]",
				ParentName:     "[[.sbomReport.metadata.namespace]]",
				ParentVersion:  "[[.sbomReport.report.artifact.tag]]",
				ProjectTags:    []string{},
			},
			input: sbomReportV1alpha1,
			mockUploadBOM: mockUploadBOM{
				enable:         true,
				projectName:    "library/alpine",
				projectVersion: "latest",
				parentName:     "default",
				parentVersion:  "latest",
				projectTags: []string{
					"test",
					"kube_namespace:default",
				},
				err: nil,
			},
			mockAddTagsToProject: mockAddTagsToProject{
				enable: false,
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		if tc.mockUploadBOM.enable {
			mockDTrack.EXPECT().
				UploadBOM(ctx, tc.mockUploadBOM.projectName, tc.mockUploadBOM.projectVersion, tc.mockUploadBOM.parentName, tc.mockUploadBOM.parentVersion, gomock.Any()).
				Return(tc.mockUploadBOM.err)
		}
		if tc.mockAddTagsToProject.enable {
			mockDTrack.EXPECT().
				AddTagsToProject(ctx, tc.mockAddTagsToProject.projectName, tc.mockAddTagsToProject.projectVersion, tc.mockAddTagsToProject.projectTags).
				Return(tc.mockAddTagsToProject.err)
		}
		if tc.mockDeleteSBOM.enableDeactivate {
			mockDTrack.EXPECT().
				DeactivateProject(ctx, tc.mockDeleteSBOM.projectName, tc.mockDeleteSBOM.projectVersion).
				Return(tc.mockUploadBOM.err)
		}
		if tc.mockDeleteSBOM.enableDelete {
			mockDTrack.EXPECT().
				DeleteProject(ctx, tc.mockDeleteSBOM.projectName, tc.mockDeleteSBOM.projectVersion).
				Return(tc.mockUploadBOM.err)
		}

		u := &Upload{
			dtrack: mockDTrack,
			config: tc.config,
		}

		if err := u.Run(ctx, tc.input); (err != nil) != tc.wantErr {
			t.Errorf("Upload.Run() error = %v, wantErr %v", err, tc.wantErr)
		}
	}
}
