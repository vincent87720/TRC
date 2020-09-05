WINDOWS := windows
LINUX := linux
DARWIN := darwin

build:
	make rmbin
	make release
	make cpbin


rmbin:
	rm -rf ./bin/* ./test/bin/*

cpbin:
	cp -r ./bin ./test


##########BUILD##########
.PHONY: windows
windows:
	mkdir -p ./bin/$(WINDOWS)
	cd ./cmd/trc;GOOS=$(WINDOWS) GOARCH=amd64 go build -o ../../bin/$(WINDOWS)/trc.exe

.PHONY: linux
linux:
	mkdir -p ./bin/$(LINUX)
	cd ./cmd/trc;GOOS=$(LINUX) GOARCH=amd64 go build -o ../../bin/$(LINUX)/trc

.PHONY: darwin
darwin:
	mkdir -p ./bin/$(DARWIN)
	cd ./cmd/trc;GOOS=$(DARWIN) GOARCH=amd64 go build -o ../../bin/$(DARWIN)/trc

.PHONY: release
release: windows linux darwin



##########RUN##########
.PHONY: runwindows
runWindows:
	cd ./test/bin/windows;./trc.exe

.PHONY: rundarwin
runDarwin:
	cd ./test/bin/darwin;./trc
run: runwindows runDarwin



# GOPATH := /d/Project/Go/src

# build_rsrc:
# 	cd $(GOPATH)/github.com/akavel/rsrc;go build 

# build_goversioninfo:
# 	go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
# 	cd $(GOPATH)/github.com/josephspurrier/goversioninfo/cmd/goversioninfo


# mksyso:
# 	cd ./cmd/trc;go generate

# gobinddata:
# 	go get -u github.com/jteeuwen/go-bindata/...
# 	cd ./assets;go-bindata -o ../cmd/trc/asset.go -pkg asset icon/...
