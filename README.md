
## examples
[examples](examples) contains examples of how to use the library.


## Command line interface
```shell
./xcli -h # 查看帮助
./xcli install -h # 查看安装帮助
# install service
sudo ./xcli install -r "-f,prod-config.yaml"  -r "--logLevel,debug"
# use environment variables
XCLI_LOG_LEVEL=debug sudo ./xcli install -r "-f,prod-config.yaml"  -r "--logLevel,debug"

sudo ./xcli uninstall # uninstall service
sudo ./xcli start # start service
sudo ./xcli stop # stop service
sudo ./xcli restart # restart service

./xcli log -h # 查看日志帮助
./xcli log # 查看日志
./xcli log -f # 查看日志并跟踪
```
