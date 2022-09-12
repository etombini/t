TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
VERSION := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || echo "v0.0.0")

# get the latest commit hash in the short form
COMMIT := $(shell git rev-parse --short HEAD)
# get the latest commit date in the form of YYYYmmdd
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")

ifneq ($(COMMIT), $(TAG_COMMIT))
	VERSION := $(VERSION)+$(DATE).$(COMMIT)
endif

build: 
	go build -o t -ldflags="-extldflags -static -s -w -X 'main.Version=${VERSION}'" . 

install:
	go install -o t -ldflags="-extldflags -static -s -w -X 'main.Version=${VERSION}'" . 
