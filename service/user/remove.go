package user

import (
	"github.com/gin-gonic/gin"
	"go_awd/dao"
	"go_awd/model"
	"go_awd/pkg/e"
	"go_awd/serializer"
)

// BanUserByID
// @Description: 管理员通过id去ban用户
// @receiver u *EmptyService
// @param c *gin.Context
// @param id int64
// @return serializer.Response
func (u *EmptyService) BanUserByID(c *gin.Context, id int64) serializer.Response {
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(id)
	if err != nil {
		return serializer.RespCode(e.InvalidWithNotExistUser, c)
	}
	if user.Role == model.AdminRole { // 无法越级ban管理员
		return serializer.RespCode(e.InvalidWithAuth, c)
	}

	user.Role = model.NoneRole // ban用户
	if err := userDao.UpdateByID(id, user); err != nil {
		return serializer.RespCode(e.InvalidWithUpdateUser, c)
	}
	return serializer.RespSuccess(e.SuccessWithUpdate, nil, c)
}
