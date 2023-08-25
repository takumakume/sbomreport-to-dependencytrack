.PHONY: mock
mock: ## Generate mocks for testing.
	mockgen -package=mock -source ./dependencytrack/dependencytrack.go -destination ./mock/dependencytrack_mock.go DependencyTrackClient
	mockgen -package=mock -source ./uploader/uploader.go -destination ./mock/uploader_mock.go Uploader

.PHONY: go-deps
go-deps:
	go install github.com/golang/mock/mockgen@v1.6.0

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: mock fmt vet ## Run tests.
	go test -v ./...

.PHONY: helm
helm:
	cd charts && helm package ../helm/sbomreport-to-dependencytrack && helm repo index .
