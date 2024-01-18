build:
	GOOS=linux GOARCH=amd64 go build -o bin/edge-deviceplugin cmd/plugin/main.go
