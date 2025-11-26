package common

type Response[T any] struct {
	ResponseCode string `json:"response_code"`
	ResponseMsg  string `json:"response_msg"`
	ResponseTime string `json:"response_time"`
	Data         T      `json:"data,omitempty"`
}
