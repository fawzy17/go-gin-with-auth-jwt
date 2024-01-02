package models

type Response struct {
	StatusCode int         `json:"status_code"`
	IsSuccess  bool        `json:"is_success"`
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
}
