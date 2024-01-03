package server

type SyncFileInfoReq struct {
	FilePath string `json:"filePath"`
	FromPeer string `json:"fromPeer"`
	Size     int64  `json:"size"`
}

type RemoveFileInfoReq struct {
	FilePath string `json:"filePath"`
	FromPeer string `json:"fromPeer"`
}
