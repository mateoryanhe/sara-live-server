package httpserver

import (
	"errors"
	"fmt"
	"xr-game-server/core/xrjson"
	"xr-game-server/core/xrlog"

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
			failResp.Message = errercode.GetMsg(errercode.SysError, GetLang(r))
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
	r.Response.Header().Set("Content-Type", contentTypeJson)
	r.Response.Write(result.resp)
	stashAPIResponseBufferWrittenAt(r)
	logAPIRequestResponseWrite(r, elapsedMs(writeStart), len(result.resp), respContent)
}

func finalizeHandlerError(r *ghttp.Request, err error) {
	if err == nil {
		return
	}
	if errercode.IsBusiness(err) {
		r.SetError(nil)
		return
	}
	logUnexpectedHandlerError(r, err)
}

func logUnexpectedHandlerError(r *ghttp.Request, err error) {
	if r == nil || err == nil {
		return
	}
	xrlog.ErrorWithErr(r.Context(), "Handler",
		fmt.Sprintf("time=%v,reqId=%v,authId=%v,method=%v,url=%v",
			gtime.Now().Time.Format("2006-01-02 15:04:05.000"),
			r.GetHeader(ReqId, ""),
			authIdFromRequest(r),
			r.Method,
			r.RequestURI,
		),
		err,
	)
}
