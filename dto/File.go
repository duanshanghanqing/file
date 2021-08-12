package dto

type File struct {
	FileId     string `json:"fileId"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Suffix     string `json:"suffix"`
	State      string `json:"state"`
	Bucket     string `json:"bucket"`
	Prefix     string `json:"prefix"`
	CreateTime int64  `json:"createTime"`
	UpdateTime int64  `json:"updateTime"`
}

type FilePage struct {
	PageBean
	Data File `json:"data"`
}
