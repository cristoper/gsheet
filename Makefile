version := v0.1.3
platforms := linux_386 linux_amd64 darwin_amd64 windows_386 windows_amd64
SHELL := /bin/bash

.PHONY: build
build:
	go build -o gsheet ./cmd/gsheet/

.PHONY: test
test:
	go test -v ./...

.PHONY: xbuild
xbuild:
	for platform in ${platforms}; do \
		pair=($${platform/_/ }) ;\
		GOARCH=$${pair[1]} GOOS=$${pair[0]} go build -ldflags '-s -w -X main.version=${version}' -o build/$$platform/gsheet ./cmd/gsheet/ ;\
		zip -j -r build/$${platform}.zip build/$${platform}/ ;\
	done
