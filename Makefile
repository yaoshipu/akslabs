SHELL := /bin/bash

build-linux:
	cd cmd/akslabs; GOOS=linux GOARCH=amd64 go build; cd ../..

build-darwin:
	cd cmd/akslabs; GOOS=darwin GOARCH=amd64 go build; cd ../..

build-windows:
	cd cmd/akslabs; GOOS=windows GOARCH=amd64 go build; cd ../..
