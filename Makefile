flags := -ldflags "-X main.versionShort=`git describe --tags --abbrev=0` \
	           -X main.versionLong=`git describe --tags`"

gen-grafana-panel-json: main.go vendor
	go build $(flags) -o gen-grafana-panel-json main.go

vendor:
	glide up

clean:
	rm gen-grafana-panel-json
