@echo off
set CGO_ENABLED=1
set GOOS=android
set GOARCH=arm64
set CC=E:\init\windows-ndk-x86_64\bin\aarch64-linux-android21-clang
set tmpPath=%~dp0
cd %tmpPath:~0,1%:
for %%I in ("%tmpPath%..\") do set "parentPath=%%~fI"
cd %parentPath%
@echo on
go build -buildmode=c-shared  -ldflags "-s -w" -o "%tmpPath%Library/Android/arm64-v8a/libSunny.so"
pause