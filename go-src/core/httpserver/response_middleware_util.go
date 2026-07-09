package httpserver

import (
	"errors"
	"xr-game-server/core/xrjson"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
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
		if sysError {
			failResp.Message = err.Error()
		} else {
			failResp.Message = errercode.GetMsg(errercode.XRCode(code), GetLang(r))
		}
		respData = failResp
	} else if wrapSuccess != nil {
		respData = wrapSuccess(res)
	} else {
		respData = res
	}
	resp := xrjson.MustMarshal(respData)
	return responseBuildResult{
		resp:     resp,
		code:     code,
		sysError: sysError,
	}
}

func writeResponse(r *ghttp.Request, wrapSuccess func(any) any) {
	writeStart := gtime.Now()
	result := buildResponseResult(r, wrapSuccess)
	r.Response.Header().Set("Content-Type", contentTypeJson)
	r.Response.Write(result.resp)
	stashAPIResponseBufferWrittenAt(r)
	logAPIRequestResponseWrite(r, elapsedMs(writeStart), len(result.resp))
}
