package xcli

import (
	"fmt"
	"github.com/kardianos/service"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"strings"
)

var DefaultStopCallback = func() error { return nil }
var (
	buildTime = "" //编译时自动填充，请不要人工赋值
	gitHash   = "" //编译时自动填充，请不要人工赋值
	goVersion = "" //编译时自动填充，请不要人工赋值

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
	defaultCommands := []*cli.Command{
		{
			Name:  "version",
			Usage: "show version",
			Action: func(cCtx *cli.Context) error {
				//flag.Parse()
				fmt.Println("build-time:", buildTime)
				fmt.Println("build-hash:", gitHash)
				fmt.Println("build-go  :", goVersion)
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
				//		Usage:       "the working directory path, must be an absolute path",
				//		Required:    true,
				//		Destination: &workingDirectory,
				//	},
				//},
				&cli.StringSliceFlag{
					// ./xxx install -r '-f,/path/of/config.yaml' # 逗号分隔key和value
					Name:    "argument",
					Aliases: []string{"r"},
					Usage:   "the arguments for install",
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
