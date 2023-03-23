package challenge

import (
	"github.com/gin-gonic/gin"
	"go_awd/dao"
	"go_awd/pkg/e"
	"go_awd/serializer"
)

func (s *UpdateChallengeInfoService) UpdateChallenge(c *gin.Context, id int64) serializer.Response {
	chalDao := dao.NewChallengeDao(c)
	chal, err := chalDao.GetByID(id)
	if err != nil {
		return serializer.RespCode(e.InvalidWithNotExistChallenge, c)
	}
	if chal.State != "success" {
		return serializer.RespCode(e.InvalidWithNotSuccessChallenge, c)
	}
	flag := false
	if s.Title != "" {
		chal.Title = s.Title
		flag = true
	}
	if s.Info != "" {
		chal.Info = s.Info
		flag = true
	}
	if s.BaseScore != 0 {
		chal.BaseScore = s.BaseScore
		flag = true
	}
	if s.Type != "" {
		chal.Type = s.Type
		flag = true
	}
	if s.InnerServerPort != "" {
		chal.InnerServerPort = s.InnerServerPort
		flag = true
	}
	if flag {
		if err := chalDao.UpdateByID(chal); err != nil {
			return serializer.RespCode(e.InvalidWithUpdateChallenge, c)
		}
	}

	return serializer.RespSuccess(e.SuccessWithUpdate, serializer.BuildChallenge(chal), c)
}
