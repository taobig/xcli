
# 以下划线命名，这样在auto-complete时不会出现
#_buildFlag = -ldflags "-s -w -X 'main.buildTime=`date -u '+%Y-%m-%dT%H:%M:%SZ'`' -X main.gitHash=`git rev-parse --short HEAD` -X 'main.goVersion=`go version`' "
_buildFlag = -ldflags "-s -w -X 'github.com/taobig/xcli.buildTime=`date -u '+%Y-%m-%dT%H:%M:%SZ'`' -X github.com/taobig/xcli.gitHash=`git rev-parse --short HEAD` -X 'github.com/taobig/xcli.goVersion=`go version`' "

all: help

.PHONY: build
build:
	@CGO_ENABLED=0 go build -o main $(_buildFlag) .
	$(if $(shell command -v upx), upx main)

.PHONY: buildlinux
buildlinux:
	@GOOS=linux CGO_ENABLED=0 go build -o main-linux $(_buildFlag) .
	@$(if $(shell command -v upx), upx main-linux)

.PHONY: buildwin
buildwin:
	@GOOS=windows CGO_ENABLED=0 go build -o main.exe $(_buildFlag) .
	@$(if $(shell command -v upx), upx main.exe)

.PHONY: builddarwin
builddarwin:
	@GOOS=darwin CGO_ENABLED=0 go build -o main-darwin $(_buildFlag) .
	@$(if $(shell command -v upx), upx main-darwin)

.PHONY: buildall
buildall: buildlinux buildwin builddarwin

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	@go clean
#	@rm -vf main 多此一举，因为go clean已经做了
#	@rm -vf main.exe 多此一举，因为go clean已经做了
	@rm -vf main-linux
	@rm -vf main-darwin

.PHONY: help
help:
	@echo "make build 编译go代码生成当前平台二进制文件"
	@echo "make buildlinux 编译go代码生成linux平台的二进制文件"
	@echo "make buildwin 编译go代码生成windows平台的二进制文件"
	@echo "make builddarwin 编译go代码生成mac平台的二进制文件"
	@echo "make buildall 编译go代码生成所有平台的二进制文件"
	@echo "make clean 清理中间目标文件"
	@echo "make test 执行测试case"

