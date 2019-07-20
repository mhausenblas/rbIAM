rbiam_version:= v0.1

.PHONY: build

build:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=${rbiam_version}" -o release/rbiam-macos .
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${rbiam_version}" -o release/rbiam-linux .
