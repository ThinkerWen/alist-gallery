package common

type FsUP struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	Path string      `json:"path"`
	Name string      `json:"name"`
	Url  string      `json:"url"`
	Task interface{} `json:"task"`
}
