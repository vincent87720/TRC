all: rmbin mksyso build cpbin

rmbin:
	rm -f bin\trc.exe test\trc.exe

mksyso:
	go generate TRC\cmd\trc

build:
	go build -o TRC\bin\trc.exe TRC\cmd\trc

cpbin:
	cp bin\trc.exe test