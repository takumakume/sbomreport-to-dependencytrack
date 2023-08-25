.PHONY: mock
mock:
	mockgen -package=mock -source ./dependencytrack/dependencytrack.go -destination ./mock/dependencytrack_mock.go DependencyTrackClient
	mockgen -package=mock -source ./uploader/uploader.go -destination ./mock/uploader_mock.go Uploader

