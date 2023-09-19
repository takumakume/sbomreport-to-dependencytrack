package dependencytrack

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"time"

	dtrack "github.com/DependencyTrack/client-go"
)

type DependencyTrackClient interface {
	UploadBOM(ctx context.Context, projectName, projectVersion string, bom []byte) error
	AddTagsToProject(ctx context.Context, projectName, projectVersion string, tags []string) error
	GetProjectForNameVersion(ctx context.Context, projectName, projectVersion string, excludeInactive, onlyRoot bool) (p dtrack.Project, err error)
}

type DependencyTrack struct {
	Client *dtrack.Client
}

var (
	ErrProjectNotFound = errors.New("project not found")
)

func IsNotFound(err error) bool {
	switch err {
	case ErrProjectNotFound:
		return true
	}

	switch err := err.(type) {
	case *dtrack.APIError:
		if err.StatusCode == 404 {
			return true
		}
	}

	return false
}

func New(baseURL, apiKey string, timeout time.Duration) (*DependencyTrack, error) {
	client, err := dtrack.NewClient(baseURL, dtrack.WithAPIKey(apiKey), dtrack.WithTimeout(timeout))
	if err != nil {
		return nil, err
	}

	return &DependencyTrack{
		Client: client,
	}, nil
}

func (dt *DependencyTrack) UploadBOM(ctx context.Context, projectName, projectVersion string, bom []byte) error {
	log.Printf("Uploading BOM: project %s:%s", projectName, projectVersion)

	uploadToken, err := dt.Client.BOM.Upload(ctx, dtrack.BOMUploadRequest{
		ProjectName:    projectName,
		ProjectVersion: projectVersion,
		AutoCreate:     true,
		BOM:            base64.StdEncoding.EncodeToString(bom),
	})
	if err != nil {
		return err
	}

	log.Printf("Polling completion of upload BOM: project %s:%s token %s", projectName, projectVersion, uploadToken)

	doneChan := make(chan struct{})
	errChan := make(chan error)

	go func(ctx context.Context) {
		defer func() {
			close(doneChan)
			close(errChan)
		}()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				processing, err := dt.Client.BOM.IsBeingProcessed(ctx, uploadToken)
				if err != nil {
					errChan <- err
					return
				}
				if !processing {
					doneChan <- struct{}{}
					return
				}
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			}
		}
	}(ctx)

	select {
	case <-doneChan:
		log.Printf("BOM upload completed: project %s:%s token %s", projectName, projectVersion, uploadToken)
		break
	case err := <-errChan:
		log.Printf("Error: BOM upload failed: project %s:%s token %s: %s", projectName, projectVersion, uploadToken, err)
		return err
	}

	return nil
}

func (d *DependencyTrack) GetProjectForNameVersion(ctx context.Context, projectName, projectVersion string, excludeInactive, onlyRoot bool) (p dtrack.Project, err error) {
	projects, err := d.Client.Project.GetProjectsForName(ctx, projectName, excludeInactive, onlyRoot)
	if err != nil {
		return p, err
	}
	for _, project := range projects {
		if project.Version == projectVersion {
			return project, nil
		}
	}
	return p, ErrProjectNotFound
}

func (dt *DependencyTrack) AddTagsToProject(ctx context.Context, projectName, projectVersion string, tags []string) error {
	log.Printf("Adding tags to project. project %s:%s tags %v", projectName, projectVersion, tags)

	project, err := dt.Client.Project.Lookup(ctx, projectName, projectVersion)
	if err != nil {
		return err
	}

	for _, tag := range tags {
		project.Tags = append(project.Tags, dtrack.Tag{Name: tag})
	}

	_, err = dt.Client.Project.Update(ctx, project)
	if err != nil {
		return err
	}

	return nil
}
