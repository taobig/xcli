package main

import (
	"fmt"
	"github.com/taobig/xcli"
	"github.com/taobig/xcli/examples/internal"
	"github.com/urfave/cli/v2"
	"os"
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
		EnvVars: map[string]string{
			"ENV_VAR": "value",
		},
	}

	startAction := func(cCtx *cli.Context) {
		fmt.Println("arg test:", cCtx.Bool("test")) //获取参数示例一
		fmt.Println("arg configFile:", configFile)  //获取参数示例二

		fmt.Println("ENV_VAR:", os.Getenv("ENV_VAR")) //获取环境变量示例

		fmt.Println("programme start")

		// 如果执行完以上内容就希望结束进程
		err := internal.KillProcessItself()
		if err != nil {
			panic(err)
		}
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
		&cli.BoolFlag{
			// ./xx --test  cCtx.Bool("test") => true
			// ./xx			cCtx.Bool("test") => false
			Name: "test",
			//Aliases:     []string{"t"},
			Usage: "test参数，不使用设置Destination方式接收",
			Value: false,
		},
	}
	var commands []*cli.Command
	xCli.Start(flags, commands)
	return
}
