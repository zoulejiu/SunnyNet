set CGO_ENABLED=1
set GOOS=android
set GOARCH=386
set CC=C:\Microsoft\AndroidNDK\android-ndk-r23c\toolchains\llvm\prebuilt\windows-x86_64\bin\i686-linux-android16-clang
go build -buildmode=c-shared  -ldflags "-s -w" -o Library/x86/libSunny.so
pause