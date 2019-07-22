EXECUTABLE=jsnschema-custom-validator
VERSION=${shell git describe --tags --always}
LDFLAGS= -s -w -X '${shell pwd}/version.Version=$(VERSION)'
BUILD_DIR=bin
.PHONY: clean dep build test all
all: dep build test

clean:
	rm -rf $(BUILD_DIR)/

dep:
	dep ensure -v --vendor-only
build:
	mkdir -p bin/
	CGO_ENABLED=0 GOOS=darwin go build -ldflags "$(LDFLAGS)" -a -o $(BUILD_DIR)/$(EXECUTABLE)
test:
	go test ./...
