GIT_TAG=$(shell git describe --tags --abbrev=0)
GIT_HASH=$(shell git rev-parse HEAD)
BUILD_DATE=$(shell date '+%F-%H:%M:%S')

info:
	@echo "[web-annie info]\nbuild information : ${GIT_TAG} - ${GIT_HASH} (${BUILD_DATE})"

build:
	@echo "[web-annie build]\nbuild information : ${GIT_TAG} - ${GIT_HASH} (${BUILD_DATE})"
	@go build -v -x -ldflags "-X main.gitTag=${GIT_TAG} -X main.gitHash=${GIT_HASH} -X main.buildDate=${BUILD_DATE}"