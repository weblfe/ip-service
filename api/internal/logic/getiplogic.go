package logic

import (
	"context"
	"crypto/md5"
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

func (l *GetIpLogic) GetIp(req types.GetIpRequest) (*types.GetIpResponse, error) {
	var region = l.getDb()
	defer l.closeDb()
	ip, err := region.MemorySearch(req.IpAddr)
	if err != nil {
		return l.sendFail(), err
	}
	return l.sendSuccess(req.IpAddr, ip), nil
}

func (l *GetIpLogic) sendSuccess(ip string, info ip2region.IpInfo) *types.GetIpResponse {
	return &types.GetIpResponse{
		Msg:  "OK",
		Code: 0,
		Data: &types.GetIpResponseData{
			Ip: ip, CityId: info.CityId,
			Country: info.Country, District: info.Region,
			Province: info.Province, City: info.City, ISP: info.ISP,
		},
	}
}

func (l *GetIpLogic) sendFail() *types.GetIpResponse {
	return &types.GetIpResponse{
		Msg:  "failed",
		Data: nil,
		Code: 1400,
	}
}

func (l *GetIpLogic) openDb(path string) (*ip2region.Ip2Region, error) {
	if p := l.readLock(path); p != "" {
		path = p
	}
	if strings.Contains(path, "http") {
		var url = path
		res, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		path, err = l.saveTemp(res.Body, time.Now().Unix())
		if err != nil {
			return nil, err
		}
		defer l.setLock(getMd5(url) + "=>" + path)
	}
	return ip2region.New(path)
}

func (l *GetIpLogic) saveTemp(closer io.ReadCloser, timestamp int64) (string, error) {
	var err error
	var tmpDir = l.svcCtx.Config.TempDir
	var file = joinPath(tmpDir, fmt.Sprintf("%d_ip2region.db", timestamp))
	defer closer.Close()
	file, err = filepath.Abs(file)
	if err != nil {
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

func (l *GetIpLogic) closeDb() {
	if l.ipDb != nil {
		l.ipDb.Close()
	}
}

func (l *GetIpLogic) isLock() bool {
	var lock = l.getLockFile()
	if _, err := os.Stat(lock); os.IsNotExist(err) {
		return false
	}
	return true
}

func (l *GetIpLogic) getLockFile() string {
	f,err:=filepath.Abs(joinPath(l.svcCtx.Config.TempDir, "db.lock"))
	if err!=nil {
		logx.Error(err)
		return ""
	}
	return f
}

// 设置存储锁
func (l *GetIpLogic) readLock(args ...string) string {
	if l.isLock() {
		data, err := ioutil.ReadFile(l.getLockFile())
		if err != nil {
			return ""
		}
		info := string(data)
		arr := strings.Split(info, "=>")
		if _, err := os.Stat(arr[1]); os.IsExist(err) {
			return ""
		}
		// url 地址不一致
		if len(args) > 0 && getMd5(args[0]) != arr[0] {
			return ""
		}
		return arr[1]
	}
	return ""
}

// 自动锁,不更新ip db
func (l *GetIpLogic) setLock(db string) bool {
	if l.isLock() {
		return false
	}
	fs, err := os.OpenFile(l.getLockFile(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logx.Error(err)
		return false
	}
	defer fs.Close()
	n, err := fs.Write([]byte(db))
	if  err == nil && n > 0 {
		return true
	}
	return false
}

// 路径处理
func joinPath(root, sub string) string {
	var file = filepath.Join(root, sub)
	if strings.Contains(file, "//") {
		file = strings.ReplaceAll(file, "//", "/")
	}
	if strings.Contains(file, "\\/") {
		file = strings.ReplaceAll(file, "\\/", string(os.PathSeparator))
	}
	return file
}

// 获取md5
func getMd5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}
