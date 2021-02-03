package logic

import (
	"context"
	"fmt"
	"github.com/mohong122/ip2region/binding/golang/ip2region"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ip-service/api/internal/svc"
	"ip-service/api/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetIpLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	ipDb   *ip2region.Ip2Region
}

func NewGetIpLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetIpLogic {
	return GetIpLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetIpLogic) getDb() *ip2region.Ip2Region {
	var err error
	if l.ipDb == nil {
		l.ipDb, err = l.openDb(l.svcCtx.Config.Database)
		if err != nil {
			l.Error(err)
		}
	}
	return l.ipDb
}

func (l *GetIpLogic) GetIp(req types.Request) (*types.Response, error) {
	var region = l.getDb()
	ip, err := region.MemorySearch(req.IpAddr)
	if err != nil {
		return l.sendFail(), err
	}
	return l.sendSuccess(req.IpAddr, ip), nil
}

func (l *GetIpLogic) sendSuccess(ip string, info ip2region.IpInfo) *types.Response {
	return &types.Response{
		Msg:  "OK",
		Code: 0,
		Data: &types.ResponseData{
			Ip: ip, CityId: info.CityId,
			Country: info.Country, District: info.Region,
			Province: info.Province, City: info.City, ISP: info.ISP,
		},
	}
}

func (l *GetIpLogic) sendFail() *types.Response {
	return &types.Response{
		Msg:  "failed",
		Data: nil,
		Code: 1400,
	}
}

func (l *GetIpLogic) openDb(path string) (*ip2region.Ip2Region, error) {
	if strings.Contains(path, "http") {
		res, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		path, err = l.saveTemp(res.Body, time.Now().Unix())
		if err != nil {
			return nil, err
		}
	}
	return ip2region.New(path)
}

func (l *GetIpLogic) saveTemp(closer io.ReadCloser, timestamp int64) (string, error) {
	var err error
	var tmpDir = l.svcCtx.Config.TempDir
	var file = fmt.Sprintf("%s/%d_ip2region.db", tmpDir, timestamp)
	if strings.Contains(file,"//") {
		file = strings.ReplaceAll(file,"//","/")
	}
	if strings.Contains(file,"\\/") {
		file = strings.ReplaceAll(file,"\\/",string(os.PathSeparator))
	}
	defer closer.Close()
	file ,err= filepath.Abs(file)
	if err!=nil {
		return "", err
	}
	dir := filepath.Dir(file)
	_, err = os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			l.Error(err)
			return "", err
		}
	}
	var fs *os.File

	var buf []byte
	var n int
	fs, err = os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return "", err
	}
	buf, err = ioutil.ReadAll(closer)
	if err != nil {
		return "", err
	}
	n, err = fs.Write(buf)
	defer fs.Close()
	if err == nil && n > 0 {
		return file, nil
	}
	return "", err
}
