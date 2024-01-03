package server

import (
	"github.com/koihuang/speedfs/config"
	"github.com/koihuang/speedfs/util"
	"net/http"
	"net/url"
	"path"
	"strings"
)

func (server *Server) Download(w http.ResponseWriter, r *http.Request) {

	var err error
	innerValue := r.FormValue("inner")
	if innerValue != "" { // download from peer
		server.staticHandler.ServeHTTP(w, r)
		return
	}
	uri := r.RequestURI
	filepath := strings.Split(uri, "?")[0]
	filepath, err = url.PathUnescape(filepath)
	if err != nil {
		writeFailRes(w, systemErr, err.Error())
		return
	}
	if !util.FileExist(path.Join(server.fileRootDir, r.RequestURI)) {
		for _, peer := range config.OtherPeers() {
			err := server.syncFromPeer(peer, FileInfo{
				Filepath: filepath,
				Size:     0,
			})
			if err != nil {
				continue
			} else {
				break
			}
		}
	}
	server.staticHandler.ServeHTTP(w, r)
}
