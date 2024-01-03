package server

import (
	"github.com/koihuang/speedfs/config"
	"github.com/koihuang/speedfs/util"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"net/http"
	"path"
	"strconv"
)

var log = logrus.StandardLogger()

type Server struct {
	fileRootDir   string
	staticHandler http.Handler
}

func New() (*Server, error) {
	fileRootDir := path.Join(config.SpeedfsPath(), config.FILE_DIR_NAME)
	return &Server{
		fileRootDir:   fileRootDir,
		staticHandler: http.StripPrefix("/", http.FileServer(http.Dir(fileRootDir))),
	}, nil
}

func (server *Server) Start() error {
	server.initRouter()

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(config.GloableConfig.Port),
		Handler: &HttpHandler{},
	}

	util.Go(server.repairyAt3AmEveryDay)

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) repairyAt3AmEveryDay() {

	c := cron.New()
	_, err := c.AddFunc("0 3 * * *", func() {
		server.repairLastDay()
	})
	if err != nil {
		log.Errorf("add cron err:%s", err.Error())
	}
	c.Start()
}
