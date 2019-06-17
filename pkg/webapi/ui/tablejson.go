package ui

type tableJson struct {
	Code  int `json:"code"`
	Count int `json:"count"`
	Data  []interface{} `json:"data"`
	Msg   string `json:"msg"`
}
