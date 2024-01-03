package server

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/google/uuid"
	"github.com/koihuang/speedfs/config"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	gopath "path"
	"strings"
	"time"
)

func (server *Server) Upload(w http.ResponseWriter, r *http.Request) {

	filepath, err := server.upload(r)
	if err != nil {
		log.Error(fmt.Sprintf("upload err: %s", err.Error()))
		writeFailRes(w, systemErr, err.Error())
	} else {
		writeSuccessRes(w, UploadRes{Filepath: filepath, DownloadUrl: fmt.Sprintf("http://%s/%s", config.SelfPeer(), filepath)})
	}

}

func (server *Server) upload(r *http.Request) (string, error) {
	var (
		err          error
		filename     string
		uploadFile   multipart.File
		uploadHeader *multipart.FileHeader
		fullpath     string
	)
	if uploadFile, uploadHeader, err = r.FormFile("file"); err != nil {
		fmt.Println(err)
	}
	_, filename = gopath.Split(uploadHeader.Filename)
	currentTime := time.Now()
	formattedTime := currentTime.Format("150405")
	uniqueID := formattedTime + "_" + strings.ReplaceAll(uuid.New().String(), "-", "")
	folder := time.Now().Format("20060102")

	var relativePath string
	if filename == "" {
		relativePath = gopath.Join(folder, uniqueID)
	} else {
		relativePath = gopath.Join(folder, uniqueID+"_"+filename)
	}

	fullpath = gopath.Join(server.fileRootDir, relativePath)
	err = os.MkdirAll(gopath.Dir(fullpath), 0775)
	if err != nil {
		return "", err
	}
	file, err := os.OpenFile(fullpath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(file, uploadFile)
	if err != nil {
		return "", err
	}

	stat, err := os.Stat(fullpath)
	if err != nil {
		return "", err
	}
	server.notifyPeersToSync(SyncFileInfoReq{
		FilePath: relativePath,
		FromPeer: config.GloableConfig.AdvertiseIP,
		Size:     stat.Size(),
	})
	return gopath.Clean("/" + relativePath), nil
}

func (server *Server) notifyPeersToSync(fileInfo SyncFileInfoReq) {
	peers := config.OtherPeers()
	data, err := json.Marshal(fileInfo)
	if err != nil {
		log.Error(err)
		return
	}
	for _, peer := range peers {
		go func(p string) {
			req := httplib.Post(fmt.Sprintf("http://%s/sync", p))
			req.Param("fileInfo", string(data))
			_, err = req.Response()
			if err != nil {
				log.Errorf("fail to sync file to peer, peer:%s, path:%s, err: %s", p, fileInfo.FilePath, err.Error())
			}
		}(peer)
	}
}
