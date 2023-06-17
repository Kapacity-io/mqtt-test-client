BINARY_DIRECTORY=./dist
BINARY_NAME=mqtt-test-client
VERSION=1.0.0
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X=main.Version=${VERSION} -X=main.Build=${BUILD}"

.PHONY: all
all: prepare clean mac linux windows

mac:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_DIRECTORY}/${BINARY_NAME}-mac-amd64 main.go

linux:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_DIRECTORY}/${BINARY_NAME}-linux-amd64 main.go

windows:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_DIRECTORY}/${BINARY_NAME}-windows-amd64.exe main.go

clean:
	@if [ -f ${BINARY_DIRECTORY}/${BINARY_NAME}-mac-amd64 ] ; then rm ${BINARY_DIRECTORY}/${BINARY_NAME}-mac-amd64 ; fi
	@if [ -f ${BINARY_DIRECTORY}/${BINARY_NAME}-linux-amd64 ] ; then rm ${BINARY_DIRECTORY}/${BINARY_NAME}-linux-amd64 ; fi
	@if [ -f ${BINARY_DIRECTORY}/${BINARY_NAME}-windows-amd64.exe ] ; then rm ${BINARY_DIRECTORY}/${BINARY_NAME}-windows-amd64.exe ; fi

prepare:
	mkdir -p ./dist
	go mod download	
