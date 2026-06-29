package httpserver

import (
	"encoding/json"
	"errors"

	"github.com/gogf/gf/v2/net/ghttp"
	"xr-game-server/errercode"
)

type responseBuildResult struct {
	resp     []byte
	code     int
	sysError bool
}

func buildResponseResult(r *ghttp.Request, wrapSuccess func(any) any) responseBuildResult {
	var (
		res   = r.GetHandlerResponse()
		err   = r.GetError()
		param any
	)
	code := errercode.Success
	sysError := false
	if err != nil {
		var gErr *errercode.XError
		errors.As(err, &gErr)
		if gErr != nil {
			code = int(gErr.Code())
			param = gErr.Param
		} else {
			r.Response.ClearBuffer()
			code = errercode.SysError
			sysError = true
		}
	}

	var respData any
	if err != nil {
		failResp := CreateFailAndParam(code, param)
		failResp.Message = errercode.GetMsg(errercode.XRCode(code), GetLang(r))
		respData = failResp
	} else if wrapSuccess != nil {
		respData = wrapSuccess(res)
	} else {
		respData = res
	}
	resp, _ := json.MarshalIndent(respData, "", " ")
	return responseBuildResult{
		resp:     resp,
		code:     code,
		sysError: sysError,
	}
}

func writeResponseAndLog(r *ghttp.Request, authId string, wrapSuccess func(any) any) {
	result := buildResponseResult(r, wrapSuccess)
	r.Response.Header().Set("Content-Type", contentTypeJson)
	r.Response.Write(result.resp)

	stashAPIRequestEndLog(r, &apiRequestEndPending{
		AuthId:   authId,
		SysError: result.sysError,
		Code:     result.code,
		Resp:     result.resp,
	})
}
