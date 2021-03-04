package handler

import (
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/thinkeridea/go-extend/exnet"
	"net"
	"net/http"
	"net/url"
	"strings"

	"ip-service/api/internal/logic"
	"ip-service/api/internal/svc"
	"ip-service/api/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

var headerArr = []string{"X-REAL-IP", "X-FORWARDED-FOR", "REMOTE-ADDR"}

func GetIpHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetIpRequest
		if !setQueryForm(r, &req) {
			if err := httpx.Parse(r, &req); svc.AssertError(err) {
				httpx.OkJson(w, errorRes(err))
				return
			}
		}

		l := logic.NewGetIpLogic(r.Context(), ctx)
		if req.IpAddr == "" {
			req.IpAddr = getRealIp(r)
		}
		if exnet.HasLocalIPddr(req.IpAddr) {
			httpx.OkJson(w, errorRes(fmt.Errorf("%s is local ip",req.IpAddr)))
			return
		}
		resp, err := l.GetIp(req)
		if svc.AssertError(err) {
			httpx.OkJson(w, errorRes(err))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

// 获取客户段ip
func getRealIp(r *http.Request) string {
	var ip string
	for _, key := range headerArr {
		ip = r.Header.Get(key)
		if ip != "" {
			break
		}
	}
	if ip != "" && strings.Contains(ip, ",") {
		ipArr := strings.Split(ip, ",")
		ip = ipArr[0]
	}
	if ip == "" {
		if host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
			return host
		}
	}
	return ip
}

// 设置Query form
func setQueryForm(r *http.Request, request *types.GetIpRequest) bool {
	var query, err = url.ParseQuery(r.URL.RawQuery)
	if svc.AssertError(err) {
		logx.Info(err)
		return false
	}
	if ip := query.Get("ipAddr"); ip != "" {
		request.IpAddr = ip
	}
	if ip := query.Get("ip"); ip != "" {
		request.IpAddr = ip
	}
	if accessKey := query.Get("accessKey"); accessKey != "" {
		request.AccessKey = accessKey
	}
	if request.IpAddr == "" {
		request.IpAddr = getRealIp(r)
	}
	return true
}

// 错误
func errorRes(err error) interface{} {
	return &types.GetIpResponse{
		Msg:  err.Error(),
		Code: 500,
		Data: nil,
	}
}
