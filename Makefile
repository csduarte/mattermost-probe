.PHONY: build build-linux build-osx build-windows install run

GO=go

build: build-windows build-osx build-linux
	
build-linux: .prebuild
	@echo Build Linux amd64 at $(GOPATH)/bin/mattermost-probe
	env GOOS=linux GOARCH=amd64 $(GO) install .

build-osx: .prebuild
	@echo Build OSX amd64 at $(GOPATH)/bin/linux_amd64/mattermost-probe
	env GOOS=darwin GOARCH=amd64 $(GO) install .

build-windows: .prebuild
	@echo Build Windows amd64 at $(GOPATH)/bin/windows_amd64/mattermost-probe.exe
	env GOOS=windows GOARCH=amd64 $(GO) install .

run: .prebuild
	@echo Building and Running
	$(GO) build
	./mattermost-probe	 

install:
	glide install

.prebuild:
	@echo Ensuring no duplication of gorilla websockets
	rm -rf vendor/github.com/mattermost/platform/vendor/github.com/gorilla/websocket