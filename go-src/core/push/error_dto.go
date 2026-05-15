package push

import "xr-game-server/errercode"

type ErrorDto struct {
	Code  int `json:"code"`
	Param any `json:"param"`
}

func NewErrorDto(err errercode.XRCode, param any) *ErrorDto {
	return &ErrorDto{
		Code:  int(err),
		Param: param,
	}
}
