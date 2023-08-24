package cmd

import (
	"time"
)

type config struct {
	baseURL        string
	apiKey         string
	projectName    string
	projectVersion string
	projectTags    []string
	timeout        time.Duration
}

func newConfig(baseURL, apiKey, projectName, projectVersion string, projectTags []string, timeout int) *config {
	return &config{
		baseURL:        baseURL,
		apiKey:         apiKey,
		projectName:    projectName,
		projectVersion: projectVersion,
		projectTags:    projectTags,
		timeout:        time.Duration(timeout) * time.Second,
	}
}
