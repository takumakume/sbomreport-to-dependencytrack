package config

import (
	"errors"
	"strings"
	"time"
)

type Config struct {
	BaseURL string
	APIKey  string

	ProjectName    string
	ProjectVersion string
	ProjectTags    []string

	DtrackClientTimeout     time.Duration
	SBOMUploadTimeout       time.Duration
	SBOMUploadCheckInterval time.Duration
}

var ErrAPIKeyIsRequired = errors.New("api-key is required")

func New(baseURL, apiKey, projectName, projectVersion string, projectTags []string, dtrackClientTimeoutSec, sbomUploadTimeoutSec, sbomUploadCheckIntervalSec float64) *Config {
	if len(projectTags) == 1 && strings.Contains(projectTags[0], ",") {
		projectTags = strings.Split(projectTags[0], ",")
	}

	return &Config{
		BaseURL:                 baseURL,
		APIKey:                  apiKey,
		ProjectName:             projectName,
		ProjectVersion:          projectVersion,
		ProjectTags:             projectTags,
		DtrackClientTimeout:     time.Duration(dtrackClientTimeoutSec) * time.Second,
		SBOMUploadTimeout:       time.Duration(sbomUploadTimeoutSec) * time.Second,
		SBOMUploadCheckInterval: time.Duration(sbomUploadCheckIntervalSec) * time.Second,
	}
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		return ErrAPIKeyIsRequired
	}

	return nil
}
