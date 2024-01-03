package cmd

import (
	"encoding/json"
	"errors"
	"github.com/koihuang/speedfs/config"
	"github.com/koihuang/speedfs/util"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

func Init(ctx *cli.Context) error {
	var (
		err        error
		configPath string
		localIP    string
		peers      []string
	)

	configDir := path.Join(config.SPEEDFS_PATH, config.CONF_DIR_NAME)
	configPath = path.Join(configDir, config.CONF_FILE_NAME)
	exist := util.FileExist(configPath)
	if exist {
		return errors.New("already be inited")
	}
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return err
	}
	localIP, err = util.LocalIP()
	if err != nil {
		localIP = ""
	}

	if localIP != "" {
		peers = append(peers, localIP+":"+strconv.Itoa(config.DEFAULT_SERVER_PORT))
	}

	initConfig := config.Config{
		Port:        config.DEFAULT_SERVER_PORT,
		AdvertiseIP: localIP,
		Peers:       peers,
	}
	cfgJson, err := json.MarshalIndent(initConfig, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configPath, cfgJson, 0644)
	if err != nil {
		return err
	}

	fileRootDir := path.Join(config.SPEEDFS_PATH, config.FILE_DIR_NAME)
	err = os.MkdirAll(fileRootDir, 0755)
	if err != nil {
		return err
	}

	logRootDir := path.Join(config.SPEEDFS_PATH, config.LOG_DIR_NAME)
	err = os.MkdirAll(logRootDir, 0755)
	if err != nil {
		return err
	}
	return nil
}
