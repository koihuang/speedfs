package config

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

var (
	BUILD_TIME string
	GIT_COMMIT string
	GO_VERSION string
)

var (
	SPEEDFS_PATH  string
	GloableConfig Config
)

const (
	SPEEDFS_PATH_ENV     = "SPEEDFS_PATH"
	SPEEDFS_PATH_DEFAULT = "~/.speedfs"
	FILE_DIR_NAME        = "file"
	CONF_DIR_NAME        = "config"
	CONF_FILE_NAME       = "config.json"
	LOG_DIR_NAME         = "log"
	PID_FILE_NAME        = "speedfs.pid"
	DEFAULT_SERVER_PORT  = 9999
)

type Config struct {
	Port        int      `json:"port"`
	AdvertiseIP string   `json:"advertise-ip"`
	Peers       []string `json:"peers"`
}

func init() {
	rootPath := os.Getenv(SPEEDFS_PATH_ENV)
	var err error
	if len(rootPath) == 0 {
		rootPath, err = homedir.Expand(SPEEDFS_PATH_DEFAULT)
	}
	if err != nil {
		panic(err)
	}

	logRotate := &lumberjack.Logger{
		Filename:   "log/speedfs.log",
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
	logrus.SetOutput(logRotate)
}

func ConfDir() string {
	return path.Join(SPEEDFS_PATH, CONF_DIR_NAME)
}

func ConfFilePath() string {
	return path.Join(SPEEDFS_PATH, ConfDir(), CONF_FILE_NAME)
}

func LogDir() string {
	return path.Join(SPEEDFS_PATH, LOG_DIR_NAME)
}

func GlobalLogger() *logrus.Logger {
	return logrus.StandardLogger()
}

func PidFilePath() string {
	return path.Join(SPEEDFS_PATH, PID_FILE_NAME)
}

func SpeedfsPath() string {
	return SPEEDFS_PATH
}

func LoadGlobalConfig() error {
	fi, err := os.Lstat(ConfFilePath())
	if fi != nil || (err != nil && !os.IsNotExist(err)) {
		cfgFile, err := os.Open(ConfFilePath())
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(cfgFile)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &GloableConfig)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("config file not found")
	}
}

func SelfPeer() string {
	return GloableConfig.AdvertiseIP + ":" + strconv.Itoa(GloableConfig.Port)
}

func OtherPeers() []string {
	var otherPeers []string
	for _, peer := range GloableConfig.Peers {
		if strings.Contains(peer, GloableConfig.AdvertiseIP+":") { // exclude self
			continue
		}
		otherPeers = append(otherPeers, peer)
	}
	return otherPeers
}

func PrintSysteminfo() {
	GlobalLogger().Info("speedfs.BUILD_TIME:", BUILD_TIME)
	GlobalLogger().Info("speedfs.GIT_COMMIT:", GIT_COMMIT)
	GlobalLogger().Info("speedfs.GO_VERSION:", GO_VERSION)
}
