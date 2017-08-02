.PHONY: build build-linux build-osx build-windows install 

GO=go

build: build-linux build-windows build-osx
	
build-linux: .prebuild
	@echo Build Linux amd64
	env GOOS=linux GOARCH=amd64 $(GO) install .

build-osx: .prebuild
	@echo Build OSX amd64
	env GOOS=darwin GOARCH=amd64 $(GO) install .

build-windows: .prebuild
	@echo Build Windows amd64
	env GOOS=windows GOARCH=amd64 $(GO)  .

install:
	glide install

.prebuild:
	@echo Ensuring no duplication of gorilla websockets
	rm -rf vendor/github.com/mattermost/platform/vendor/github.com/gorilla/websocket