package cmd

import (
	"time"

	"github.com/takumakume/sbomreport-to-dependencytrack/dependencytrack"
)

type config struct {
	dtrack         *dependencytrack.DependencyTrack
	projectName    string
	projectVersion string
	projectTags    []string
}

func newConfig(baseURL, apiKey, projectName, projectVersion string, projectTags []string, timeout int) (*config, error) {
	dtrack, err := dependencytrack.New(baseURL, apiKey, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	return &config{
		dtrack:         dtrack,
		projectName:    projectName,
		projectVersion: projectVersion,
		projectTags:    projectTags,
	}, nil
}
