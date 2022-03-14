SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o .\bin\darwin\trc .\cmd\trc

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o .\bin\linux\trc .\cmd\trc

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o .\bin\windows\trc.exe .\cmd\trc
cd .\cmd\TRCGUI\ && go generate && cd ..\..\
cd .\cmd\TRCGUI\ && go build -ldflags="-H windowsgui" -o ..\..\bin\windows\TRCGUI.exe && cd ..\..\