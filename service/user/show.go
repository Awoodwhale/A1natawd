package user

import (
	"github.com/gin-gonic/gin"
	"go_awd/dao"
	"go_awd/pkg/e"
	wjwt "go_awd/pkg/wjwt"
	"go_awd/serializer"
)

func (u *EmptyService) ShowByID(c *gin.Context, id int64) serializer.Response {
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(id)
	if err != nil {
		return serializer.RespCode(e.InvalidWithShow, c)
	}
	return serializer.RespSuccess(e.SuccessWithShow, serializer.BuildUser(user), c)
}

func (u *EmptyService) ShowSelf(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	return u.ShowByID(c, claims.ID)
}