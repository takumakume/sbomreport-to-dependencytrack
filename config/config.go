package config

import (
	"errors"
	"strings"
)

type Config struct {
	BaseURL string
	APIKey  string

	ProjectName    string
	ProjectVersion string
	ProjectTags    []string
}

var ErrAPIKeyIsRequired = errors.New("api-key is required")

func New(baseURL, apiKey, projectName, projectVersion string, projectTags []string) *Config {
	if len(projectTags) == 1 && strings.Contains(projectTags[0], ",") {
		projectTags = strings.Split(projectTags[0], ",")
	}

	return &Config{
		BaseURL:        baseURL,
		APIKey:         apiKey,
		ProjectName:    projectName,
		ProjectVersion: projectVersion,
		ProjectTags:    projectTags,
	}
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		return ErrAPIKeyIsRequired
	}

	return nil
}
