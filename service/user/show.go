package user

import (
	"github.com/gin-gonic/gin"
	"go_awd/conf"
	"go_awd/dao"
	"go_awd/pkg/e"
	wjwt "go_awd/pkg/wjwt"
	"go_awd/serializer"
	"go_awd/service"
)

// ShowByID
// @Description: 通过id获取用户信息
// @receiver u *EmptyService
// @param c *gin.Context
// @param id int64
// @return serializer.Response
func (u *EmptyService) ShowByID(c *gin.Context, id int64) serializer.Response {
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(id)
	if err != nil {
		return serializer.RespCode(e.InvalidWithShow, c)
	}
	return serializer.RespSuccess(e.SuccessWithShow, serializer.BuildUser(user), c)
}

// ShowSelf
// @Description: 获取当前登录用户的信息
// @receiver u *EmptyService
// @param c *gin.Context
// @return serializer.Response
func (u *EmptyService) ShowSelf(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	return u.ShowByID(c, claims.ID)
}

// ShowUsers
// @Description: 管理员的操作，获取用户列表
// @receiver u *EmptyService
// @param c *gin.Context
// @return serializer.Response
func (u *EmptyService) ShowUsers(c *gin.Context) serializer.Response {
	if u.PageSize == 0 {
		u.PageSize = conf.PageSize
	}
	userDao := dao.NewUserDao(c)
	users, err := userDao.ListByCondition(nil, &u.BasePage)
	if err != nil {
		service.Errorln("ShowUsers ListByCondition error,", err.Error())
		return serializer.RespCode(e.InvalidWithShow, c)
	}
	return serializer.RespList(e.SuccessWithShow, serializer.BuildUsers(users), len(users), c)
}
