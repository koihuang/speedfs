package server

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/koihuang/speedfs/config"
	"net/http"
	"os"
	"path"
)

func (server *Server) Delete(w http.ResponseWriter, r *http.Request) {
	err := server.delete(r)
	if err != nil {
		log.Error(fmt.Sprintf("delete file err: %s", err.Error()))
		writeFailRes(w, systemErr, err.Error())
	} else {
		writeSuccessRes(w, nil)
	}
}

func (server *Server) delete(r *http.Request) error {
	filepath := r.FormValue("filepath")
	defer server.notifyPeersToDelete(filepath)
	fullpath := path.Join(server.fileRootDir, filepath)
	err := os.Remove(fullpath)
	if err != nil {
		return nil
	}
	return nil
}

func (server *Server) notifyPeersToDelete(filepath string) {
	peers := config.OtherPeers()

	for _, peer := range peers {
		go func(p string) {
			req := httplib.Post(fmt.Sprintf("http://%s/remove", p))
			req.Param("filepath", filepath)
			_, err := req.Response()
			if err != nil {
				log.Errorf("fail to notify peer to delete file, peer:%s, filepath:%s, err: %s", p, filepath, err.Error())
			}
		}(peer)
	}
}
