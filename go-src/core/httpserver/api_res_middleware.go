package httpserver

import (
	"encoding/json"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"strconv"
	"time"
	"xr-game-server/constants/common"
	"xr-game-server/errercode"
)

func apiResponseMiddleware(r *ghttp.Request) {
	r.Middleware.Next()
	doTime := r.GetHeader(DoTime)
	if len(doTime) == common.Zero {
		doTime = strconv.FormatInt(time.Now().UnixMilli(), 10)
	}
	retTime := time.Now().UnixMilli() - gconv.Int64(doTime)
	var (
		res   = r.GetHandlerResponse()
		err   = r.GetError()
		ctx   = gctx.New()
		param any
	)
	code := errercode.Success
	if err != nil {
		var gErr *errercode.XError
		errors.As(err, &gErr)
		if gErr != nil {
			code = int(gErr.Code())
			param = gErr.Param
			g.Log().Infof(ctx, "错误码应答,处理时间=%vms,url=%v,ip=%v,authId=%v，响应错误码数据=%v", retTime, r.URL.RequestURI(), r.Host, r.GetHeader(AuthId), code)
		} else {
			//清除系统自带的错误日志输出
			r.Response.ClearBuffer()
			code = errercode.SysError
			g.Log().Infof(ctx, "出现无法捕获的错误,处理时间=%vms,url=%v,ip=%v,authId=%v,请求数据=%v", retTime, r.URL.RequestURI(), r.Host, r.GetHeader(AuthId), requestBodyForLog(r))
		}
	}
	r.Response.Header().Set("Content-Type", contentTypeJson)
	var respData any
	//发现有错误码
	if err != nil {
		failResp := CreateFailAndParam(code, param)
		failResp.Message = errercode.GetMsg(errercode.XRCode(code), GetLang(r))
		respData = failResp
	} else {
		respData = CreateSuccess(res)
	}
	resp, _ := json.MarshalIndent(respData, "", " ")
	r.Response.Write(resp)
	if err == nil {
		if retTime >= LongDoTime {
			g.Log().Infof(ctx, "请求处理时间过长,处理时间=%vms,url=%v,ip=%v,authId=%v,响应数据=%v", retTime, r.URL.RequestURI(), r.Host, r.GetHeader(AuthId), responseBodyForLog(r, resp))
		} else {
			g.Log().Infof(ctx, "正常响应,处理时间=%vms,url=%v,ip=%v,authId=%v，响应数据=%v", retTime, r.URL.RequestURI(), r.Host, r.GetHeader(AuthId), responseBodyForLog(r, resp))
		}
	}

}
