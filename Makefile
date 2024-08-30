.PHONY : build clean build_android_arm64 build_linux_amd64 build_windows_amd64 build_darwin_amd64

build: build_android_arm64 build_linux_amd64 build_windows_amd64 build_darwin_amd64
	@echo All Done

build_android_arm64:
	GOOS=android GOARCH=arm64 go build -o aha-android-arm64

build_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o aha-linux-amd64

build_windows_amd64:
	GOOS=windows GOARCH=amd64 go build -o aha-windows-amd64.exe

build_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -o aha-darwin-amd64

clean:
	@rm -rf aha-*
