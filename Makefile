flags := -ldflags "-X main.versionShort=`git describe --tags --abbrev=0` \
	           -X main.versionLong=`git describe --tags`"

gen-grafana-panel-json: main.go vendor
	go build $(flags) -o gen-grafana-panel-json main.go

vendor:
	glide up

clean:
	rm gen-grafana-panel-json

version := `./gen-grafana-panel-json -version`
release:
	ghr -u tmtk75 --prerelease $(version) pkg

XC_ARCH := amd64
XC_OS := linux darwin windows
build:
	for arch in $(XC_ARCH); do \
	  for os in $(XC_OS); do \
	    echo $$arch $$os; \
	    GOARCH=$$arch GOOS=$$os go build \
	      -o pkg/gen-grafana-panel-json_$${os}_$$arch \
	      $(flags) main.go \
	  done \
	done
