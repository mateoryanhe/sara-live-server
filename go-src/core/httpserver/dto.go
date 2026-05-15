package httpserver

import "xr-game-server/errercode"

type Response struct {
	Message string      `json:"message" dc:"消息提示"`
	Data    interface{} `json:"data"    dc:"执行结果"`
	Code    int         `json:"code"    dc:"错误码"`
}

func CreateFail(code int) *Response {
	return &Response{
		Code: code,
	}
}

func CreateFailAndParam(code int, param any) *Response {
	return &Response{
		Code: code,
		Data: param,
	}
}

func CreateSuccess(data interface{}) *Response {
	return &Response{
		Data: data,
		Code: errercode.Success,
	}
}

func CreateSuccessNonData() *Response {
	return &Response{
		Code: errercode.Success,
	}
}

type CMSQueryReq struct {
	PageIndex int `json:"pageIndex" dc:"页码"`
	PageSize  int `json:"pageSize"  dc:"页大小"`
}

type CMSQueryResp struct {
	Total int         `json:"total" dc:"总数"`
	Data  interface{} `json:"data"  dc:"列表"`
}

func NewCMSQueryResp(total int, data interface{}) *CMSQueryResp {
	return &CMSQueryResp{
		Total: total,
		Data:  data,
	}
}
