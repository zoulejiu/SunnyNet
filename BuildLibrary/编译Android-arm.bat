set CGO_ENABLED=1
set GOOS=android
set GOARCH=arm
set CC=C:\Microsoft\AndroidNDK\android-ndk-r23c\toolchains\llvm\prebuilt\windows-x86_64\bin\armv7a-linux-androideabi21-clang
go build -buildmode=c-shared  -ldflags "-s -w" -o "Library/armeabi-v7a/libSunny.so"
pause