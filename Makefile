EXECUTABLE=jsnschema-custom-validator
VERSION=${shell git describe --tags --always}
PACKAGE=github.com/igor-karpukhin/jsnschema-custom-validator
LDFLAGS= -s -w -X '$(PACKAGE)/version.Version=$(VERSION)'
BUILD_DIR=bin
.PHONY: clean dep build test all
all: dep build test

clean:
	rm -rf $(BUILD_DIR)/

dep:
	dep ensure -v --vendor-only
build:
	mkdir -p bin/
	for target in "darwin" "linux" "windows" ; do \
		CGO_ENABLED=0 GOOS=$$target go build -ldflags "$(LDFLAGS)" -a -o $(BUILD_DIR)/$(EXECUTABLE)_$$target ; \
	done
test:
	go test ./...
