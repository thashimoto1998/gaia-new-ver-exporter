build:
	GOOS=darwin GOARCH=arm64 go build -o new-ver-exportor-darwin-arm64
	GOOS=linux GOARCH=arm64 go build -o new-ver-exportor-linux-arm64
	GOOS=linux GOARCH=amd64 go build -o new-ver-exportor-linux-amd64