package svc

import (
	"github.com/tal-tech/go-zero/core/logx"
	"ip-service/api/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}

func AssertError(err error) bool  {
	if err==nil {
		return false
	}
	logx.Error(err)
	return true
}