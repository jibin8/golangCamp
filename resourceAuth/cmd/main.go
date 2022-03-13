package main

import (
	"context"
	"github.com/urfave/cli/v2"
	"os"
	"self/internal"
	"self/pkg/logger"
)

var VERSION = "1.1.1"

func main() {
	ctx := logger.CreateTraceIDContext(context.Background(), "resource-auth-main")

	appInstance := cli.NewApp()
	appInstance.Name = "ResourceAuth"
	appInstance.Version = VERSION
	appInstance.Commands = []*cli.Command{
		newCmd(ctx),
	}
	err := appInstance.Run(os.Args)
	if err != nil {
		logger.Errorf(ctx, err.Error())
	}
}

func newCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "resourceAuth",
		Usage: "启动服务端",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "配置文件(.json,.yaml,.toml)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "model",
				Aliases:  []string{"m"},
				Usage:    "权限配置(.conf)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "menu",
				Usage: "初始化菜单数据配置文件(.yaml)",
			},
		},
		Action: func(c *cli.Context) error {
			return internal.Run(ctx,
				internal.SetConfigFile(c.String("conf")),
				internal.SetModelFile(c.String("model")),
				internal.SetMenuFile(c.String("menu")),
				internal.SetVersion(VERSION))
		},
	}
}
