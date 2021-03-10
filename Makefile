WINDOWS := windows
LINUX := linux
DARWIN := darwin
PACKAGE := github.com/vincent87720/TRC/cmd/trc

build: rmbin release cpbin

rmbin:
	rm -rf ./bin/* ./test/bin/*

cpbin:
	cp -r ./bin ./test


##########BUILD##########
.PHONY: buildwindows
buildwindows:
	GOOS=$(WINDOWS) GOARCH=amd64 go build -o bin/$(WINDOWS)/trc.exe $(PACKAGE)

.PHONY: buildlinux
buildlinux:
	GOOS=$(LINUX) GOARCH=amd64 go build -o bin/$(LINUX)/trc $(PACKAGE)

.PHONY: builddarwin
builddarwin:
	GOOS=$(DARWIN) GOARCH=amd64 go build -o bin/$(DARWIN)/trc $(PACKAGE)

.PHONY: release
release: buildwindows buildlinux builddarwin



##########RUN##########
run:
	go run $(PACKAGE)

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
