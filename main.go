package main

import (
	"fmt"
	"github.com/koihuang/speedfs/cmd"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "speedfs"
	app.Commands = []*cli.Command{
		{
			Name:   "init",
			Usage:  "init speedfs",
			Action: cmd.Init,
		},
		{
			Name:   "start",
			Usage:  "start speedfs",
			Action: cmd.Start,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
