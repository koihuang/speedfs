package server

import (
	"github.com/astaxie/beego/httplib"
	"github.com/koihuang/speedfs/config"
	"path"
	"time"
)

type ListFileByDirRes struct {
	JsonResult
	Data []*FileInfo `json:"data"`
}

func (server *Server) repairLastDay() {
	log.Info("start to repair last day's files")
	currentTime := time.Now()
	yesterday := currentTime.AddDate(0, 0, -1).Format("20060102")
	fullDir := path.Join(server.fileRootDir, yesterday)
	yesterdayFiles, err := server.listFullDir(fullDir, yesterday)
	if err != nil {
		log.Errorf("fail to list self node's file, fulldir:%s err:%s", fullDir, err.Error())
	}
	yesterdayFilesMap := make(map[string]*FileInfo)
	for _, file := range yesterdayFiles {
		yesterdayFilesMap[file.Filepath] = file
	}

	for _, peer := range config.OtherPeers() {
		postUrl := "http://" + peer + "/listFileByDir"
		req := httplib.Post(postUrl)
		req.Param("dir", yesterday)
		var res ListFileByDirRes
		err = req.ToJSON(&res)
		if err != nil {
			log.Errorf("req peer:%s for listFileByDir err:%s", peer, err.Error())
		} else {
			for _, fileInfo := range res.Data {
				if _, exist := yesterdayFilesMap[fileInfo.Filepath]; !exist {
					err = server.syncFromPeer(peer, *fileInfo)
					if err != nil {
						log.Errorf("repair file from peer:%s, err:%s", peer, err.Error())
					} else {
						yesterdayFilesMap[fileInfo.Filepath] = fileInfo
					}
				}
			}
		}
	}
	log.Info("repair last day's files end")

}
