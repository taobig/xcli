
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

## linux-systemd only
./xcli log -h # Print a help message
./xcli log # (equal to `journalctl -u xxxx.service`)
./xcli log -f # (equal to `journalctl -f -u xxxx.service`)

./xcli status -h # Print a help message
./xcli status # (equal to `systemctl status xxxx.service`)
./xcli status --no-pager # (equal to `systemctl status -l --no-pager xxxx.service`)
```
