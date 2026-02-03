package common

type Response[T any] struct {
	ResponseId   string `json:"response_id"`
	ResponseCode string `json:"response_code"`
	ResponseMsg  string `json:"response_msg"`
	ResponseTime string `json:"response_time"`
	Data         T      `json:"data"`
}
