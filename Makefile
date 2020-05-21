all: rmbin mksyso build cpbin

rmbin:
	rm -f bin\trc.exe test\trc.exe

mksyso:
	tools\rsrc\rsrc.exe -manifest tools\rsrc\main.manifest -ico=assets\icon\blockchain-blueblue.ico -o cmd\trc\rsrc.syso

build:
	go build -o .\bin\trc.exe TRC\cmd\trc

cpbin:
	cp bin\trc.exe test