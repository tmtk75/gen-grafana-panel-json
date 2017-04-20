cmd := gen-grafana-panel-json

VERSION := $(shell git describe --tags --abbrev=0)
VERSION_LONG := $(shell git describe --tags)
flags := -ldflags "-X main.versionShort=$(VERSION) \
	           -X main.versionLong=$(VERSION_LONG)"

$(cmd): *.go vendor
	go build $(flags) -o $(cmd)  *.go

install: *.go vendor
	go install $(flags)

vendor:
	glide up

clean:
	rm $(cmd)

release: pkg gen-grafana-panel-json
	ghr -u tmtk75 --prerelease $(VERSION) pkg

XC_ARCH := amd64
XC_OS := linux darwin windows
pkg: *.go
	rm -f pkg/*.gz
	for arch in $(XC_ARCH); do \
	  for os in $(XC_OS); do \
	    echo $$arch $$os; \
	    GOARCH=$$arch GOOS=$$os go build \
	      -o pkg/$(cmd)_$${os}_$$arch \
	      $(flags) `\ls *.go | grep -v _test.go`; \
	  done; \
	done
	gzip pkg/*
	touch pkg
