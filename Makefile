# FIXME: it would be nice to encode branch too
VERSION  := $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)

all: mint

mint: *.go cmd/*.go go.*
	go build -o mint -ldflags "-X github.com/colinnewell/jenkins-job-mint/cmd.Version=$(VERSION)" main.go

test:
	go test ./...

install: mint
	cp mint /usr/local/bin

lint:
	golangci-lint run
	gofmt -l -s .

license-check:
	# gem install license_finder
	license_finder
