package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func (server *Server) ListFileByDir(w http.ResponseWriter, r *http.Request) {
	fileInfos, err := server.listDir(r)
	if err != nil {
		log.Error(fmt.Sprintf("upload err: %s", err.Error()))
		writeFailRes(w, systemErr, err.Error())
	} else {
		writeSuccessRes(w, fileInfos)
	}
}

func (server *Server) listDir(r *http.Request) ([]*FileInfo, error) {
	dir := r.FormValue("dir")
	if dir == "" {
		return nil, errors.New("dir param should not be empty")
	}
	fullDir := path.Join(server.fileRootDir, dir)
	return server.listFullDir(fullDir, dir)
}

func (server *Server) listFullDir(fullDir, userDir string) ([]*FileInfo, error) {
	var fileInfos []*FileInfo
	_, err := os.Stat(fullDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fileInfos, nil
		} else {
			return nil, err
		}
	}
	fileInfoList, err := ioutil.ReadDir(fullDir)
	if err != nil {
		return nil, err
	}
	for _, info := range fileInfoList {
		if info.IsDir() {
			continue
		}
		fileInfos = append(fileInfos, &FileInfo{
			Filepath: path.Join("/", userDir, info.Name()),
			Size:     info.Size(),
		})
	}
	return fileInfos, nil
}
