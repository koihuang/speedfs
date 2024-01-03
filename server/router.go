package server

import "net/http"

var mux = new(http.ServeMux)

func (server *Server) initRouter() {
	mux.HandleFunc("/upload", server.Upload)
	mux.HandleFunc("/delete", server.Delete)
	mux.HandleFunc("/", server.Download)
	mux.HandleFunc("/listFileByDir", server.ListFileByDir)
	mux.HandleFunc("/sync", server.Sync)
	mux.HandleFunc("/remove", server.Remove)
}
