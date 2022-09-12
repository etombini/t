TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
# `2>/dev/null` suppress errors and `|| true` suppress the error codes.
TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || echo "v0.0.0")
# here we strip the version prefix
VERSION := $(TAG)

# get the latest commit hash in the short form
COMMIT := $(shell git rev-parse --short HEAD)
# get the latest commit date in the form of YYYYmmdd
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")

# check if the version string is empty
ifeq ($(VERSION),)
    VERSION := $(DATE)-$(COMMIT)
endif

ifneq ($(COMMIT), $(TAG_COMMIT))
    VERSION := $(VERSION)-$(DATE)-$(COMMIT)
endif


build: 
	go build -ldflags="-extldflags -static -s -w -X 'main.Version=${VERSION}'" . 

install:
	go install -ldflags="-extldflags -static -s -w -X 'main.Version=${VERSION}'" . 
