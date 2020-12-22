NAME=psutils
PROJECT_ID=psutils
IMAGE_TAG=gcr.io/$(PROJECT_ID)/$(NAME)

dev:
	realize start --name=$(NAME)

dependencies:
	go mod download
	go mod tidy

update-dependencies:
	go get -u ./...
	go mod tidy

test: dependencies
	go test ./...

build: dependencies
	go build main.go

debug:
	dlv debug --headless --listen=:2345 --log --api-version 2

release:
	goreleaser --rm-dist

release-dev:
	goreleaser --rm-dist --snapshot
