package uploader

import (
	"context"
	"errors"
	"log"

	"github.com/takumakume/sbomreport-to-dependencytrack/config"
	"github.com/takumakume/sbomreport-to-dependencytrack/dependencytrack"
	"github.com/takumakume/sbomreport-to-dependencytrack/sbomreport"
	tmpl "github.com/takumakume/sbomreport-to-dependencytrack/template"
)

type Uploader interface {
	Run(ctx context.Context, input []byte) error
}

type Upload struct {
	dtrack dependencytrack.DependencyTrackClient
	config *config.Config
}

func New(c *config.Config) (*Upload, error) {
	dtrack, err := dependencytrack.New(c.BaseURL, c.APIKey, c.DtrackClientTimeout, c.SBOMUploadTimeout, c.SBOMUploadCheckInterval)
	if err != nil {
		return nil, err
	}

	return &Upload{
		dtrack: dtrack,
		config: c,
	}, nil
}

func (u *Upload) Run(ctx context.Context, input []byte) error {
	sbom, err := sbomreport.New(input)
	if err != nil {
		if sbomreport.IsErrNotSBOMReport(err) {
			log.Printf("SKIP: %s", err)
			return nil
		}
		return err
	}

	if !sbom.ISVerbUpdate() {
		return errors.New("only support verb is update")
	}

	sbomMap, err := sbom.ToMap()
	if err != nil {
		return err
	}

	tpl := tmpl.New(sbomMap)

	projectName, err := tpl.Render(u.config.ProjectName)
	if err != nil {
		return err
	}

	projectVersion, err := tpl.Render(u.config.ProjectVersion)
	if err != nil {
		return err
	}

	projectTags := []string{}
	for _, tag := range u.config.ProjectTags {
		t, err := tpl.Render(tag)
		if err != nil {
			return err
		}
		projectTags = append(projectTags, t)
	}

	parentName, err := tpl.Render(u.config.ParentName)
	if err != nil {
		return err
	}

	parentVersion, err := tpl.Render(u.config.ParentVersion)
	if err != nil {
		return err
	}

	if err := u.dtrack.UploadBOM(ctx, projectName, projectVersion, parentName, parentVersion, sbom.BOM()); err != nil {
		return err
	}

	if len(projectTags) > 0 {
		if err := u.dtrack.AddTagsToProject(ctx, projectName, projectVersion, projectTags); err != nil {
			return err
		}
	}

	return nil
}
