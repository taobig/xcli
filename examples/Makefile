
# 以下划线命名，这样在auto-complete时不会出现
_targetName = main
#_buildFlag = -ldflags "-s -w -X 'main.buildTime=`date -u '+%Y-%m-%dT%H:%M:%SZ'`' -X main.gitHash=`git rev-parse --short HEAD` "
_buildFlag = -ldflags "-s -w -X 'github.com/taobig/xcli.buildTime=`date -u '+%Y-%m-%dT%H:%M:%SZ'`' -X github.com/taobig/xcli.gitHash=`git rev-parse --short HEAD` "
_buildTag = -tags urfave_cli_no_docs # https://cli.urfave.org/#urfave_cli_no_docs

all: help

.PHONY: build
build:
	@CGO_ENABLED=0 go build --trimpath -o $(_targetName) $(_buildFlag) $(_buildTag) .
#	$(if $(shell command -v upx), upx $(_targetName))

.PHONY: buildlinux
buildlinux:
	@GOOS=linux CGO_ENABLED=0 go build --trimpath -o $(_targetName)-linux $(_buildFlag) $(_buildTag) .
	@$(if $(shell command -v upx), upx $(_targetName)-linux)

.PHONY: buildwin
buildwin:
	@GOOS=windows CGO_ENABLED=0 go build --trimpath -o $(_targetName).exe $(_buildFlag) $(_buildTag) .
	@$(if $(shell command -v upx), upx $(_targetName).exe)

.PHONY: builddarwin
builddarwin:
	@GOOS=darwin CGO_ENABLED=0 go build --trimpath -o $(_targetName)-darwin $(_buildFlag) $(_buildTag) .
#	@$(if $(shell command -v upx), upx $(_targetName)-darwin)

.PHONY: buildall
buildall: buildlinux buildwin builddarwin

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	@go clean # -x The -x flag causes clean to print remove commands as it executes them.
	@rm -f $(_targetName) $(_targetName)-linux $(_targetName).exe $(_targetName)-darwin

.PHONY: help
help:
	@echo "make build 编译go代码生成当前平台二进制文件"
	@echo "make buildlinux 编译go代码生成linux平台的二进制文件"
	@echo "make buildwin 编译go代码生成windows平台的二进制文件"
	@echo "make builddarwin 编译go代码生成mac平台的二进制文件"
	@echo "make buildall 编译go代码生成所有平台的二进制文件"
	@echo "make clean 清理中间目标文件"
	@echo "make test 执行测试case"

