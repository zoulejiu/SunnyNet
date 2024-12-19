@echo off
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=386
set tmpPath=%~dp0
cd %tmpPath:~0,1%:
for %%I in ("%tmpPath%..\") do set "parentPath=%%~fI"
cd %parentPath%
@echo on
go build -ldflags "-s -w" -o "%tmpPath%Library\windows\x32\__Sunny.exe"