package server

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

func (server *Server) Remove(w http.ResponseWriter, r *http.Request) {
	err := server.remove(r)
	if err != nil {
		log.Error(fmt.Sprintf("delete file err: %s", err.Error()))
		writeFailRes(w, systemErr, err.Error())
	} else {
		writeSuccessRes(w, nil)
	}
}

func (server *Server) remove(r *http.Request) error {
	filepath := r.FormValue("filepath")
	fullpath := path.Join(server.fileRootDir, filepath)
	err := os.Remove(fullpath)
	if err != nil {
		return nil
	}
	return nil
}
