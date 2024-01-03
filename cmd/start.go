package cmd

import (
	"fmt"
	"github.com/koihuang/speedfs/config"
	"github.com/koihuang/speedfs/server"
	"github.com/koihuang/speedfs/util"
	"github.com/urfave/cli/v2"
)

func Start(ctx *cli.Context) error {

	err := config.LoadGlobalConfig()
	if err != nil {
		return err
	}

	fmt.Println("speedfs start...")
	config.PrintSysteminfo()

	err = util.WritePid(config.PidFilePath())
	if err != nil {
		return err
	}

	s, err := server.New()
	if err != nil {
		return err
	}
	err = s.Start()
	if err != nil {
		return err
	}
	return nil
}
