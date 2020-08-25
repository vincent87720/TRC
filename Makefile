GOPATH := /d/Project/Go/src

run: gorun
build: rmbin gobuild cpbin

build_rsrc:
	cd $(GOPATH)/github.com/akavel/rsrc;go build 

build_goversioninfo:
	go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
	cd $(GOPATH)/github.com/josephspurrier/goversioninfo/cmd/goversioninfo

rmbin:
	rm -f ./bin/trc.exe ./test/trc.exe

mksyso:
	cd ./cmd/trc;go generate

gobinddata:
	go get -u github.com/jteeuwen/go-bindata/...
	cd ./assets;go-bindata -o ../cmd/trc/asset.go -pkg asset icon/...

gobuild:
	cd ./cmd/trc;go build -o ../../bin/trc.exe 

cpbin:
	cp ./bin/trc.exe test

gorun:
	cd ./test;./trc.exe