package server

import (
	"encoding/json"
	"github.com/astaxie/beego/httplib"
	"github.com/koihuang/speedfs/config"
	"github.com/koihuang/speedfs/util"
	"net/http"
	"path"
	"strings"
	"time"
)

func (server *Server) Sync(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.RemoteAddr, "127.0.0.1") { // do not download file from local req
		writeSuccessRes(w, nil)
		return
	}
	err := server.sync(r)
	if err != nil {
		log.Errorf("sync file error:%s", err)
		writeFailRes(w, systemErr, err.Error())
	} else {
		writeSuccessRes(w, nil)
	}
}

func (server *Server) sync(r *http.Request) error {
	fileInfoJsonStr := r.FormValue("fileInfo")
	var syncFileInfo SyncFileInfoReq
	err := json.Unmarshal([]byte(fileInfoJsonStr), &syncFileInfo)
	if err != nil {
		return err
	}
	for _, peer := range config.OtherPeers() {
		if strings.Contains(peer, syncFileInfo.FromPeer) {
			err = server.syncFromPeer(peer, FileInfo{
				Filepath: syncFileInfo.FilePath,
				Size:     syncFileInfo.Size,
			})
			return err
		}
	}
	return nil
}

func (server *Server) syncFromPeer(peer string, fileInfo FileInfo) error {
	filepath := fileInfo.Filepath
	downloadUrl := "http://" + peer + path.Clean("/"+filepath)
	fullpath := path.Join(server.fileRootDir, fileInfo.Filepath)
	exist := util.FileExist(fullpath)
	if !exist {
		retryTimes := 2
		timeout := fileInfo.Size/1024/1024/1 + 30
		var err error
		for i := 0; i < retryTimes; i++ {
			req := httplib.Get(downloadUrl)
			req.Param("inner", "1")
			if fileInfo.Size > 0 {
				req.SetTimeout(time.Second*30, time.Second*time.Duration(timeout))
			}
			var response *http.Response
			response, err = req.Response()
			if response.StatusCode == 200 {
				err = req.ToFile(fullpath)
				if err != nil {
					log.Errorf("fail to write file err:%s", err.Error())
				} else {
					break
				}
			} else {
				continue
			}
		}
		return err
	}
	return nil
}
