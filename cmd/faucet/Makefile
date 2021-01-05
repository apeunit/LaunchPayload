GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)
GIT_DESCR = $(shell git describe --tags --always)
# build output folder
OUTPUTFOLDER = dist
# build paramters
OS = linux
ARCH = amd64

.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs


default: build

build: build-dist

build-dist: $(GOFILES)
	@echo build binary to $(OUTPUTFOLDER)
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static" -X main.Version=$(GIT_DESCR)' -o $(OUTPUTFOLDER)/faucet .
	@echo done

install: build
	@echo installing to $(GOPATH)/bin
	cp dist/* $(GOPATH)/bin
	@echo done

test: test-all

test-all:
	@go test $(GOPACKAGES) -v -race -coverprofile=cover.out -covermode=atomic


clean:
	@echo remove $(OUTPUTFOLDER) folder
	rm -rf $(OUTPUTFOLDER)
	@echo done
