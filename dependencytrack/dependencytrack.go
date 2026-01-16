package dependencytrack

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	dtrack "github.com/DependencyTrack/client-go"
)

type DependencyTrackClient interface {
	UploadBOM(ctx context.Context, projectName, projectVersion string, parentName string, parentVersion string, bom []byte) error
	AddTagsToProject(ctx context.Context, projectName, projectVersion string, tags []string) error
	DeactivateProject(ctx context.Context, projectName, projectVersion string) error
	DeleteProject(ctx context.Context, projectName, projectVersion string) error
}

type DependencyTrack struct {
	Client *dtrack.Client

	SBOMUploadTimeout       time.Duration
	SBOMUploadCheckInterval time.Duration
}

func New(baseURL, apiKey string, dtrackClientTimeout, sbomUploadTimeout, sbomUploadCheckInterval time.Duration) (*DependencyTrack, error) {
	client, err := dtrack.NewClient(baseURL, dtrack.WithAPIKey(apiKey), dtrack.WithTimeout(dtrackClientTimeout))
	if err != nil {
		return nil, err
	}

	return &DependencyTrack{
		Client:                  client,
		SBOMUploadTimeout:       sbomUploadTimeout,
		SBOMUploadCheckInterval: sbomUploadCheckInterval,
	}, nil
}

func (dt *DependencyTrack) UploadBOM(ctx context.Context, projectName, projectVersion string, parentName string, parentVersion string, bom []byte) error {
	log.Printf("Uploading BOM: project %s:%s", projectName, projectVersion)

	uploadToken, err := dt.Client.BOM.Upload(ctx, dtrack.BOMUploadRequest{
		ProjectName:    projectName,
		ProjectVersion: projectVersion,
		ParentName:     parentName,
		ParentVersion:  parentVersion,
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

		ticker := time.NewTicker(dt.SBOMUploadCheckInterval)
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
			case <-time.After(dt.SBOMUploadTimeout):
				errChan <- fmt.Errorf("timeout exceeded")
				return
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

func (dt *DependencyTrack) DeactivateProject(ctx context.Context, projectName, projectVersion string) error {
	log.Printf("Deactivating project. project %s:%s", projectName, projectVersion)

	project, err := dt.Client.Project.Lookup(ctx, projectName, projectVersion)
	if err != nil {
		return err
	}

	project.Active = false

	_, err = dt.Client.Project.Update(ctx, project)
	if err != nil {
		return err
	}

	return nil
}

func (dt *DependencyTrack) DeleteProject(ctx context.Context, projectName, projectVersion string) error {
	log.Printf("Deleting project. project %s:%s", projectName, projectVersion)

	project, err := dt.Client.Project.Lookup(ctx, projectName, projectVersion)
	if err != nil {
		return err
	}

	return dt.Client.Project.Delete(ctx, project.UUID)
}
