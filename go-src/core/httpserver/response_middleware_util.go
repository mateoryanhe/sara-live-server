package httpserver

import (
	"errors"
	"xr-game-server/core/xrjson"
	"xr-game-server/errercode"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

type responseBuildResult struct {
	resp     []byte
	code     int
	sysError bool
}

const handlerResponseCtxKey = "httpserver.handlerResponse"

func SetHandlerResponseData(r *ghttp.Request, data any) {
	if r != nil {
		r.SetCtxVar(handlerResponseCtxKey, data)
	}
}

func getHandlerResponseData(r *ghttp.Request) any {
	if r == nil {
		return nil
	}
	if v := r.GetCtxVar(handlerResponseCtxKey); !v.IsNil() {
		return v.Interface()
	}
	return r.GetHandlerResponse()
}

func buildResponseResult(r *ghttp.Request, wrapSuccess func(any) any) responseBuildResult {
	var (
		res   = getHandlerResponseData(r)
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
			failResp.Message = err.Error() // errercode.GetMsg(errercode.SysError, GetLang(r))
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
	handlerErr := r.GetError()
	writeStart := gtime.Now()
	result := buildResponseResult(r, wrapSuccess)
	finalizeHandlerError(r, handlerErr)
	respContent := string(result.resp)
	if isH5CryptoRequest(r) {
		writeH5ResponseBody(r, result.resp, writeStart)
		return
	}
	r.Response.Header().Set("Content-Type", contentTypeJson)
	r.Response.Write(result.resp)
	stashAPIResponseBufferWrittenAt(r)
	if !shouldSkipAPILogChain(r) {
		logAPIRequestResponseWrite(r, elapsedMs(writeStart), len(result.resp), respContent)
	}
}

func finalizeHandlerError(r *ghttp.Request, err error) {
	if err == nil {
		return
	}
	if errercode.IsBusiness(err) {
		r.SetError(nil)
	}
}
