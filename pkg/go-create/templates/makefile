NAME={{.Name}}
PROJECT_ID=psutils
IMAGE_TAG=gcr.io/$(PROJECT_ID)/$(NAME)

dev:
	realize start --name=$(NAME)

dependencies:
	go get .

test: dependencies
	go test ./...

build: dependencies
	go build main.go

debug:
	dlv debug --headless --listen=:2345 --log --api-version 2
