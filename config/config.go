package config

import (
	"errors"
	"slices"
	"strings"
	"time"
)

type Config struct {
	BaseURL string
	APIKey  string

	ProjectName    string
	ProjectVersion string
	ProjectTags    []string

	ParentName    string `json:"parentName,omitempty"`
	ParentVersion string `json:"parentVersion,omitempty"`

	DtrackClientTimeout     time.Duration
	SBOMUploadTimeout       time.Duration
	SBOMUploadCheckInterval time.Duration
	SBOMDeleteAction        string
}

var (
	ErrAPIKeyIsRequired        = errors.New("api-key is required")
	ErrInvalidSBOMDeleteAction = errors.New("invalid sbom-delete-action")
)

func New(
	baseURL, apiKey, projectName, projectVersion string,
	projectTags []string,
	parentName string,
	parentVersion string,
	dtrackClientTimeoutSec, sbomUploadTimeoutSec, sbomUploadCheckIntervalSec float64,
	sbomDeleteAction string,
) *Config {
	if len(projectTags) == 1 && strings.Contains(projectTags[0], ",") {
		projectTags = strings.Split(projectTags[0], ",")
	}

	return &Config{
		BaseURL:                 baseURL,
		APIKey:                  apiKey,
		ProjectName:             projectName,
		ProjectVersion:          projectVersion,
		ProjectTags:             projectTags,
		ParentName:              parentName,
		ParentVersion:           parentVersion,
		DtrackClientTimeout:     time.Duration(dtrackClientTimeoutSec) * time.Second,
		SBOMUploadTimeout:       time.Duration(sbomUploadTimeoutSec) * time.Second,
		SBOMUploadCheckInterval: time.Duration(sbomUploadCheckIntervalSec) * time.Second,
		SBOMDeleteAction:        sbomDeleteAction,
	}
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		return ErrAPIKeyIsRequired
	}

	if !slices.Contains([]string{"ignore", "delete", "deactivate"}, c.SBOMDeleteAction) {
		return ErrInvalidSBOMDeleteAction
	}

	return nil
}
