GO_VERSION=`go version`
GIT_COMMIT=`git rev-parse --short HEAD`
BUILD_TIME=`date "+%Y-%m-%d/%H:%M:%S"`
APP_NAME=speedfs

CONFIG_DIR = github.com/koihuang/speedfs/config

.PHONY: build pack


TAR_NAME=speedfs.tar.gz
ifeq ($(target),)
	COMPILE_TARGET=
else
	TAR_NAME=speedfs-${target}.tar.gz
    PARAMS=$(subst -, ,$(target))
    ifeq ($(words $(PARAMS)),2)
    	OS=$(word 1, $(PARAMS))
    	ARCH=$(word 2, $(PARAMS))
    	COMPILE_TARGET=CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH)
    else
        $(error error param: '$(target)'! example: 'target=darwin-amd64')
    endif
endif

build:
	${COMPILE_TARGET} go build -trimpath -buildmode=pie -ldflags "-w -s -X '${CONFIG_DIR}.BUILD_TIME=${BUILD_TIME}' -X '${CONFIG_DIR}.GO_VERSION=${GO_VERSION}' -X '${CONFIG_DIR}.GIT_COMMIT=${GIT_COMMIT}'" -o ./${APP_NAME} main.go

pack:
	make build
	rm -rf ./pack
	mkdir ./pack
	cp ./speedfs  ./pack/
	cp ./scripts/* ./pack/
	tar -czvf  ${TAR_NAME} -C ./pack/ .
	rm -rf ./pack

pack-all:
	make pack target=linux-amd64
	make pack target=linux-arm64
	make pack target=darwin-amd64
	make pack target=darwin-arm64

