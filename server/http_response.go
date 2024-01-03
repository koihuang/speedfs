package server

import (
	"encoding/json"
	"net/http"
)

type JsonResult struct {
	Status string      `json:"status"`
	Code   int         `json:"errCode"`
	Msg    string      `json:"errMsg"`
	Data   interface{} `json:"data"`
}

type UploadRes struct {
	Filepath    string `json:"filepath"`
	DownloadUrl string `json:"downloadUrl"`
}

type FileInfo struct {
	Filepath string `json:"filepath"`
	Size     int64  `json:"size"`
}

func writeFailRes(w http.ResponseWriter, code int, msg string) {
	data, err := json.Marshal(JsonResult{Status: "fail", Code: code, Msg: msg})
	if err != nil {
		log.Errorf("marshal fail res err:%s", err.Error())
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Errorf("write fail res to peer error, err:%s", err.Error())
	}
}

func writeSuccessRes(w http.ResponseWriter, data interface{}) {
	bytes, err := json.Marshal(JsonResult{Status: "success", Data: data})
	if err != nil {
		log.Errorf("marshal fail res err:%s", err.Error())
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		log.Errorf("write success res to peer error, err:%s", err.Error())
		return
	}
}
