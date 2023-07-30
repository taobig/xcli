package main

import (
	"fmt"
	"github.com/taobig/xcli"
	"github.com/urfave/cli/v2"
)

var (
	configFile string //参数变量：配置文件路径
	logLevel   string //参数变量：覆盖配置文件中的日志级别
)

func main() {
	var serviceConf = &xcli.ServiceConfig{
		ServiceName:        "demoService",
		ServiceDisplayName: "demoService",
		ServiceDescription: "...",
	}

	startAction := func() {
		fmt.Println("programme start")
	}
	stopAction := func() error {
		fmt.Println("programme exit")
		return nil
	}
	//xCli, err := xcli.New(serviceConf, startAction, xcli.DefaultStopCallback)
	xCli, err := xcli.New(serviceConf, startAction, stopAction)
	if err != nil {
		panic(err)
	}
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "configFile",
			Aliases:     []string{"f"},
			Usage:       "the config file",
			Value:       "config.yaml",
			Destination: &configFile,
		},
		&cli.StringFlag{
			Name:        "logLevel",
			Aliases:     []string{"l"},
			Usage:       "the log level, overwrite the config file's log level. e.g., debug, info, warn, error, fatal, panic",
			Value:       "",
			Destination: &logLevel,
		},
	}
	var commands []*cli.Command
	xCli.Start(flags, commands)
	return
}
