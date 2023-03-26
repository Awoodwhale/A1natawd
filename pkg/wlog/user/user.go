package user

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	c "go_awd/pkg/wlog/common"
)

func GetOpt() *c.OptLog {
	opt := c.InitOpt()
	opt.FileNamePrefix = "trace.log"
	return opt
}

func New(log *c.OptLog) (*logrus.Logger, error) {
	if lg, err := log.ConfigLogrus(); err != nil {
		return lg, errors.Cause(err)
	} else {
		return lg, nil
	}
}
