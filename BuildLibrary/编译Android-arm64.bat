set CGO_ENABLED=1
set GOOS=android
set GOARCH=arm64
set CC=C:\Microsoft\AndroidNDK\android-ndk-r23c\toolchains\llvm\prebuilt\windows-x86_64\bin\aarch64-linux-android21-clang
go build -buildmode=c-shared  -ldflags "-s -w" -o Library/arm64-v8a/libSunny.so
pause