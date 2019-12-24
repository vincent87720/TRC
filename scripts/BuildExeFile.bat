D:\Project\Go\src\github.com\akavel\rsrc\rsrc.exe -manifest main.manifest -ico=..\assest\icon\blockchain-blueblue.ico -o ..\cmd\trc\rsrc.syso
cd ..\cmd\trc
go build
copy /Y trc.exe ..\..\test\
move trc.exe ../../bin
