package handler

import (
	"github.com/tal-tech/go-zero/rest/httpx"
	"ip-service/api/internal/logic"
	"ip-service/api/internal/svc"
	"net/http"
)

func GetIndexHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetIndexLogic(r.Context(), ctx)
		resp, err := l.GetIndex()
		if err != nil {
			httpx.OkJson(w, errorRes(err))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}