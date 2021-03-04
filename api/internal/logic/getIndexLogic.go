package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/tal-tech/go-zero/core/logx"
	"io/ioutil"
	"ip-service/api/internal/svc"
)

type GetIndexLogic struct {
	Logger logx.Logger
	ctx   context.Context
	svcCtx *svc.ServiceContext
}

func NewGetIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetIndexLogic {
	return GetIndexLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIndexLogic)GetIndex()  (map[string]interface{}, error)  {
	var (
		err error
		content []byte
		decoder *json.Decoder
		docs map[string]interface{}
	)
	content,err = ioutil.ReadFile(l.svcCtx.Config.DocsPath)
	if svc.AssertError(err) {
		return nil, err
	}

	decoder = json.NewDecoder(bytes.NewBuffer(content))
	if err=decoder.Decode(&docs);svc.AssertError(err) {
		return nil,err
	}
	return docs,nil
}

