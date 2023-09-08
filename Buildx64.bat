set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64
go build -buildmode=c-shared  -ldflags "-s -w" -o "export\Sunny64.dll"
pause