package challenge

import (
	"github.com/gin-gonic/gin"
	"go_awd/conf"
	"go_awd/dao"
	"go_awd/pkg/e"
	"go_awd/serializer"
	"go_awd/service"
)

// ShowChallenges
// @Description: 获取题目列表
// @receiver s *EmptyService
// @param c *gin.Context
// @return serializer.Response
func (s *EmptyService) ShowChallenges(c *gin.Context) serializer.Response {
	if s.PageSize == 0 {
		s.PageSize = conf.PageSize
	}
	chalDao := dao.NewChallengeDao(c)
	chals, err := chalDao.ListByCondition(nil, &s.BasePage)
	if err != nil {
		service.Errorln("ShowChallenges ListByCondition error,", err.Error())
		return serializer.RespCode(e.InvalidWithShow, c)
	}
	return serializer.RespList(e.SuccessWithShow, serializer.BuildChallenges(chals), len(chals), c)
}
