package xcli

import (
	"bufio"
	"fmt"
	"github.com/kardianos/service"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var DefaultStopCallback = func() error { return nil }
var (
	buildTime = "" //编译时自动填充，请不要人工赋值
	gitHash   = "" //编译时自动填充，请不要人工赋值

	//workingDirectory string //参数变量
	installArguments cli.StringSlice //参数变量
)

type XCli struct {
	serviceConfig *ServiceConfig
	startCallback func(cCtx *cli.Context)
	stopCallback  func() error
}

func New(conf *ServiceConfig, startCallback func(cCtx *cli.Context), stopCallback func() error) (*XCli, error) {
	xCli := &XCli{
		serviceConfig: conf,
		startCallback: startCallback,
		stopCallback:  stopCallback,
	}
	return xCli, nil
}

func (x *XCli) createSystemService(cCtx *cli.Context) (service.Service, error) {
	conf := x.serviceConfig
	startCallback := x.startCallback
	stopCallback := x.stopCallback

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}
	// 默认install service时，WorkingDirectory就是其二进制程序所在的目录。这样便于加载配置文件。
	workingDir := strings.Replace(dir, "\\", "/", -1)
	logrus.Debugf("workingDir: %s", workingDir)

	logrus.Debugf("installArguments: %+v", installArguments.String())
	var arguments []string = installArguments.Value()

	svcConfig := &service.Config{
		Name:        conf.ServiceName,
		DisplayName: conf.ServiceDisplayName,
		Description: conf.ServiceDescription,
		Arguments:   arguments,
		//Executable:       "/usr/local/bin/myapp",
		//Dependencies:     []string{"After=network.target syslog.target"},
		//WorkingDirectory: "",
		WorkingDirectory: workingDir,
		//Option: service.KeyValue{
		//	"Restart": "always", // Restart=always
		//},
		EnvVars: conf.EnvVars,
	}

	ss := &SystemService{
		cCtx:          cCtx,
		startCallback: startCallback,
		stopCallback:  stopCallback,
	}
	s, err := service.New(ss, svcConfig)
	if err != nil {
		return nil, fmt.Errorf("service New failed, err: %v", err)
	}
	return s, nil
}

func (x *XCli) Start(flags []cli.Flag, commands []*cli.Command) {
	logCommandUsageText := ""
	platform := service.Platform()
	if platform == "linux-systemd" {
		logCommandUsageText = "" +
			"ServiceType: " + platform + "\n" +
			"journalctl commands:" + "\n" +
			"journalctl -u " + x.serviceConfig.ServiceName + "\n" +
			"journalctl -n 10 -u " + x.serviceConfig.ServiceName + "\n" +
			"journalctl -f -u " + x.serviceConfig.ServiceName + "\n" +
			"journalctl --disk-usage 检查当前journal使用磁盘量\n" +
			"journalctl --vacuum-time=2d 只保留两天的日志\n" +
			"journalctl --vacuum-size=1G 只保留1G的日志\n" +
			"journalctl --vacuum-files=2 只保留最近的两个日志文件\n"

	} else if platform == "darwin-launchd" {
		//logCommandUsageText = "..."
	}

	defaultCommands := []*cli.Command{
		{
			Name:  "version",
			Usage: "Show version",
			Action: func(cCtx *cli.Context) error {
				//flag.Parse()
				fmt.Println("build-time:", buildTime)
				fmt.Println("build-hash:", gitHash)
				//go1.21有了toolchain指令后，环境变量里go不一定是最终编译使用的go版本，所以请自行使用 go version -m <program> 查看实际编译使用的go版本
				fmt.Println("build-go  :", "go version -m <program>")
				return nil
			},
		},
		{
			Name:   "install",
			Action: x.controlAction,
			Flags: []cli.Flag{
				//	&cli.StringFlag{
				//		Name:        "workingDirectory",
				//		Aliases:     []string{"w"},
				//		Usage:       "The working directory path, must be an absolute path",
				//		Required:    true,
				//		Destination: &workingDirectory,
				//	},
				//},
				&cli.StringSliceFlag{
					// ./xxx install -r '-f,/path/of/config.yaml' # 逗号分隔key和value
					Name:    "argument",
					Aliases: []string{"r"},
					Usage:   "The arguments for install",
					//Value:       "",
					Destination: &installArguments,
				},
			},
		},
		{
			Name:   "uninstall",
			Action: x.controlAction,
		},
		{
			Name:   "start",
			Action: x.controlAction,
		},
		{
			Name:   "restart",
			Action: x.controlAction,
		},
		{
			Name:   "stop",
			Action: x.controlAction,
		},
		{
			Name:      "log",
			Usage:     "Show logs",
			UsageText: logCommandUsageText,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "follow",
					Aliases: []string{"f"},
					Usage:   "Show only the most recent journal entries, and continuously print new entries as they are appended to the journal.",
					Value:   false,
					//Destination: ,
				},
			},
			Action: func(cCtx *cli.Context) error {
				if platform == "linux-systemd" {
					serviceName := x.serviceConfig.ServiceName
					//journalctl -u xxxxx.service
					//journalctl -f -u xxxxx.service
					follow := cCtx.Bool("follow")
					if follow {
						cmd := exec.Command("journalctl", "-f", "-u", serviceName)
						return runCommand(cmd, "StdoutPipe")
					} else {
						cmd := exec.Command("journalctl", "-u", serviceName)
						return runCommand(cmd, "Stdout")
					}
				} else {
					fmt.Println(service.Platform() + " not support")
				}
				return nil
			},
		},
		{
			Name:  "status",
			Usage: "Show terse runtime status information about one or more units, followed by most recent log data from the journal.",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name: "no-pager",
					//Aliases: []string{"?"},
					Usage: "Do not pipe output into a pager.",
					Value: false,
					//Destination: ,
				},
			},
			Action: func(cCtx *cli.Context) error {
				if platform == "linux-systemd" {
					serviceName := x.serviceConfig.ServiceName
					//systemctl -l --no-pager status xxxxx.service
					//systemctl status xxxxx.service
					noPager := cCtx.Bool("no-pager")
					if noPager {
						cmd := exec.Command("systemctl", "-l", "--no-pager", "status", serviceName)
						return runCommand(cmd, "Stdout")
					} else {
						cmd := exec.Command("systemctl", "status", serviceName)
						return runCommand(cmd, "StdoutPipe")
					}
				} else {
					fmt.Println(service.Platform() + " not support")
				}
				return nil
			},
		},
	}
	allCommands := append(defaultCommands, commands...)
	app := &cli.App{
		Name: x.serviceConfig.ServiceName,
		//Usage: "",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			//在参数解析之后createSystemService，因为其内部可能用到了command参数
			systemService, err := x.createSystemService(cCtx)
			if err != nil {
				return err
			}

			err = systemService.Run()
			if err != nil {
				logrus.Errorf("service Run failed, err: %v", err)
				return err
			}
			return nil
		},
		Commands: allCommands,
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Errorf("Error: %+v", err)
		os.Exit(1)
	}
}

func (x *XCli) controlAction(cCtx *cli.Context) error {
	//在参数解析之后createSystemService，因为其内部可能用到了command参数
	systemService, err := x.createSystemService(cCtx)
	if err != nil {
		return err
	}
	err = service.Control(systemService, cCtx.Command.Name)
	if err != nil {
		logrus.Errorf("service %s failed, err: %v", cCtx.Command.Name, err)
		return err
	}
	return nil
}

func runCommand(cmd *exec.Cmd, outputType string) error {
	fmt.Printf("cmd: %s\n", cmd)

	if outputType == "StdoutPipe" {
		pipe, err := cmd.StdoutPipe()
		if err != nil {
			logrus.Errorf("cmd.StdoutPipe() failed with %s", err)
			return err
		}
		if err = cmd.Start(); err != nil {
			logrus.Errorf("cmd.Start() failed with %s", err)
			return err
		}
		go func(p io.ReadCloser) {
			reader := bufio.NewReader(pipe)
			line, err := reader.ReadString('\n')
			for err == nil {
				fmt.Println(line)
				line, err = reader.ReadString('\n')
			}
		}(pipe)
		if err = cmd.Wait(); err != nil {
			logrus.Errorf("cmd.Wait() failed with %s", err)
			return err
		}
	} else if outputType == "Stdout" {
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		//err := cmd.Run()
		//if err != nil {
		//	logrus.Errorf("cmd.Run() failed with %s", err)
		//	return err
		//}
		out, err := cmd.Output()
		if err != nil {
			logrus.Errorf("cmd.Output() failed with %s", err)
			return err
		}
		fmt.Println(string(out))
		//out, err := cmd.CombinedOutput()
		//if err != nil {
		//	logrus.Fatalf("cmd.Run() failed with %s", err)
		//	return err
		//}
		//fmt.Println(string(out))
	}
	return nil
}
