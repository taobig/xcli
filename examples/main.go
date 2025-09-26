package main

import (
	"context"
	"examples/internal"
	"fmt"
	"log/slog"
	"os"

	"github.com/taobig/xcli"
	"github.com/urfave/cli/v3"
)

var (
	configFile string //参数变量：配置文件路径
	logLevel   string //参数变量：覆盖配置文件中的日志级别
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	var serviceConf = &xcli.ServiceConfig{
		ServiceName:        "demoService",
		ServiceDisplayName: "demoService",
		ServiceDescription: "...",
		EnvVars: map[string]string{
			"ENV_VAR": "value",
		},
	}

	startAction := func(ctx context.Context, cmd *cli.Command) {
		fmt.Println("arg test:", cmd.Bool("test")) //获取参数示例一
		fmt.Println("arg configFile:", configFile) //获取参数示例二

		fmt.Println("ENV_VAR:", os.Getenv("ENV_VAR")) //获取环境变量示例

		fmt.Println("programme start")

		{ // optional: 如果执行完以上内容就希望结束进程，可以调用KillProcessItself。这样还是会正常触发stopAction回调的。
			err := internal.KillProcessItself()
			if err != nil {
				panic(err)
			}
		}

		// 注意：一定不要在这里执行无限循环、http ListenAndServe()等阻塞操作，否则将会导致程序在执行<start>子命令时无法正常返回

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
}
