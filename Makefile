EXECUTABLE=jsnschema-custom-validator
VERSION=${shell git describe --tags --always}
LDFLAGS= -s -w -X '${shell pwd}/version.Version=$(VERSION)'

.PHONY: dep build test all
all: dep build test

dep:
	dep ensure -v --vendor-only
build:
	CGO_ENABLED=0 GOOS=darwin go build -ldflags "$(LDFLAGS)" -a -o $(EXECUTABLE)
test:
	go test ./...
