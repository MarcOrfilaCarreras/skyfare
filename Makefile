BINARY_NAME=skyfare

.PHONY: all clean build-armv6 build-armv7 build-arm64

all: build-armv6 build-armv7 build-arm64

build-armv6:
	cd skyfare && CC=arm-linux-gnueabihf-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=6 go build -o $(BINARY_NAME)-armv6

build-armv7:
	cd skyfare && CC=arm-linux-gnueabihf-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 go build -o $(BINARY_NAME)-armv7

build-arm64:
	cd skyfare && CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME)-arm64

clean:
	cd skyfare && rm -f $(BINARY_NAME)-armv6 $(BINARY_NAME)-armv7 $(BINARY_NAME)-arm64
