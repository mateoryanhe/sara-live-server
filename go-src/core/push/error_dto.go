package push

import "xr-game-server/errercode"

type ErrorDto struct {
	Code  int `json:"code" dc:"错误码"`
	Param any `json:"param" dc:"附加参数"`
}

func NewErrorDto(err errercode.XRCode, param any) *ErrorDto {
	return &ErrorDto{
		Code:  int(err),
		Param: param,
	}
}
