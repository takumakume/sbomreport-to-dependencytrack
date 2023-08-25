package config

import "errors"

type Config struct {
	BaseURL string
	APIKey  string

	ProjectName    string
	ProjectVersion string
	ProjectTags    []string
}

var ErrAPIKeyIsRequired = errors.New("api-key is required")

func New(baseURL, apiKey, projectName, projectVersion string, projectTags []string) *Config {
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
