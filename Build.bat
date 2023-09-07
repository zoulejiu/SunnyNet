set CGO_ENABLED=1
set GOOS=windows
set GOARCH=386
go build -buildmode=c-shared  -ldflags "-s -w" -o "export\Sunny.dll"
pause