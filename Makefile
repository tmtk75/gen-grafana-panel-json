gen-grafana-panel-json: main.go vendor
	go build -o gen-grafana-panel-json main.go

vendor:
	glide up
