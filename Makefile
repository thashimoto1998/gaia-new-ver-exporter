build:
	GOOS=darwin GOARCH=arm64 go build -o new-ver-exporter-darwin-arm64
	GOOS=linux GOARCH=arm64 go build -o new-ver-exporter-linux-arm64
	GOOS=linux GOARCH=amd64 go build -o new-ver-exporter-linux-amd64