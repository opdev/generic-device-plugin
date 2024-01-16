build:
	GOOS=linux GOARCH=amd64 go build -o bin/edge-deviceplugin cmd/deviceplugin/main.go
