cmd := gen-grafana-panel-json

flags := -ldflags "-X main.versionShort=`git describe --tags --abbrev=0` \
	           -X main.versionLong=`git describe --tags`"

$(cmd): *.go vendor
	go build $(flags) -o $(cmd)  *.go

vendor:
	glide up

clean:
	rm $(cmd)

version := `./$(cmd) -version`
release: pkg gen-grafana-panel-json
	ghr -u tmtk75 --prerelease $(version) pkg

XC_ARCH := amd64
XC_OS := linux darwin windows
pkg: main.go
	rm -f pkg/*.gz
	for arch in $(XC_ARCH); do \
	  for os in $(XC_OS); do \
	    echo $$arch $$os; \
	    GOARCH=$$arch GOOS=$$os go build \
	      -o pkg/$(cmd)_$${os}_$$arch \
	      $(flags) main.go; \
	  done; \
	done
	gzip pkg/*
	touch pkg
