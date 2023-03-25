package team

import (
	"github.com/gin-gonic/gin"
	"go_awd/dao"
	"go_awd/model"
	"go_awd/pkg/e"
	wjwt "go_awd/pkg/wjwt"
	"go_awd/serializer"
	"go_awd/service"
)

// ApplyTeam
// @Description: 申请加入团队
// @receiver t *ApplyService
// @param c *gin.Context
// @return serializer.Response
func (t *ApplyService) ApplyTeam(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(claims.ID)
	if err != nil || user.TeamID != 0 || user.IsTeamLeader {
		// 没有战队才能下一步
		return serializer.RespCode(e.InvalidWithAuth, c)
	}
	// 创建或更新ut数据
	utDao := dao.NewUTDaoByDB(userDao.DB)

	if _, exist := utDao.ExistItemByUserID(user.ID); exist {
		// 已经存在战队了（审核中）
		return serializer.RespCode(e.InvalidWithReviewTeam, c)
	}

	if err := utDao.CreateOrUpdateItemByUserID(&model.UserTeam{
		UserID: user.ID,
		TeamID: t.TeamID,
		State:  0, // 审核中
	}); err != nil {
		service.Errorln("apply team mysql error, ", err.Error())
		return serializer.RespCode(e.InvalidWithApplyTeam, c)
	}

	return serializer.RespSuccess(e.SuccessWithApplyTeam, nil, c)
}

// CancelApplyTeam
// @Description: 取消战队入队申请
// @receiver t *EmptyService
// @param c *gin.Context
// @return serializer.Response
func (t *EmptyService) CancelApplyTeam(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(claims.ID)
	if err != nil || user.TeamID != 0 || user.IsTeamLeader {
		// 找不到用户、有战队了 都不能进行
		return serializer.RespCode(e.InvalidWithAuth, c)
	}

	// 查询ut信息
	utDao := dao.NewUTDaoByDB(userDao.DB)
	ut, exist := utDao.ExistItemByUserID(user.ID)
	if !exist || ut.State != 0 {
		// 已经入队过了
		return serializer.RespCode(e.InvalidWithCancelApplyTeam, c)
	}

	// 删除ut信息
	if err := utDao.DeleteByUserID(user.ID); err != nil {
		return serializer.RespCode(e.InvalidWithCancelApplyTeam, c)
	}

	return serializer.RespSuccess(e.SuccessWithCancelApplyTeam, nil, c)
}

// AcceptTeam
// @Description: 通过某个user的申请入队，只有是队长才能
// @receiver t *AcceptOrRejectService
// @param c *gin.Context
// @return serializer.Response
func (t *AcceptOrRejectService) AcceptTeam(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(claims.ID)
	if err != nil || user.TeamID == 0 || !user.IsTeamLeader {
		// 不是队长或者没有战队就不行
		return serializer.RespCode(e.InvalidWithAuth, c)
	}

	// 查询ut信息
	utDao := dao.NewUTDaoByDB(userDao.DB)
	ut, exist := utDao.ExistItemByUserID(t.UserID)
	if !exist || ut.State != 0 || ut.TeamID != user.TeamID {
		// 不存在ut或者已经在队伍中或者申请入队的不是如自己的队，就无法通过
		return serializer.RespCode(e.InvalidWithAcceptTeam, c)
	}

	// 存在就改变这个user的入队状态
	ut.State = 1 // 成为队员
	if err := utDao.CreateOrUpdateItemByUserID(ut); err != nil {
		service.Errorln("team leader accept team mysql_user_team error,", err.Error())
		return serializer.RespCode(e.InvalidWithAcceptTeam, c)
	}

	// 修改这个申请入队用户的user属性
	user, err = userDao.GetUserByID(t.UserID)
	user.TeamID = ut.TeamID
	user.IsTeamLeader = false
	if err := userDao.UpdateByID(user.ID, user); err != nil {
		service.Errorln("team leader accept team mysql_user error,", err.Error())
		return serializer.RespCode(e.InvalidWithAcceptTeam, c)
	}

	// 发送邮箱给user，告诉团队申请通过了
	go func() {
		teamDao := dao.NewTeamDaoByDB(userDao.DB)
		team, err := teamDao.GetByID(ut.TeamID)
		if err != nil {
			service.Errorln("team not found, ", err.Error())
			return
		}
		if err := service.SendUserAcceptTeam(user.Email, team.TeamName); err != nil {
			service.Errorln("SendUserAcceptTeam error, ", err.Error())
			return
		}
	}()

	return serializer.RespSuccess(e.SuccessWithAcceptTeam, nil, c)
}

// RejectTeam
// @Description: 拒绝加入团队
// @receiver t *AcceptOrRejectService
// @param c *gin.Context
// @return serializer.Response
func (t *AcceptOrRejectService) RejectTeam(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(claims.ID)
	if err != nil || user.TeamID == 0 || !user.IsTeamLeader {
		// 不是队长或者没有战队就不行
		return serializer.RespCode(e.InvalidWithAuth, c)
	}

	// 查询ut信息
	utDao := dao.NewUTDaoByDB(userDao.DB)
	ut, exist := utDao.ExistItemByUserID(t.UserID)
	if !exist || ut.State != 0 || ut.TeamID != user.TeamID {
		// 不存在ut或者已经在队伍中或者申请入队的不是如自己的队，就无法通过
		return serializer.RespCode(e.InvalidWithRejectTeam, c)
	}

	// 删除ut信息
	if err := utDao.DeleteByUserID(user.ID); err != nil {
		return serializer.RespCode(e.InvalidWithRejectTeam, c)
	}

	// 发送邮箱给user，告诉团队申请被拒绝了
	go func() {
		teamDao := dao.NewTeamDaoByDB(userDao.DB)
		team, err := teamDao.GetByID(ut.TeamID)
		if err != nil {
			service.Errorln("team not found, ", err.Error())
			return
		}
		if err := service.SendUserRejectTeam(user.Email, team.TeamName); err != nil {
			service.Errorln("SendUserRejectTeam error, ", err.Error())
			return
		}
	}()

	return serializer.RespSuccess(e.SuccessWithRejectTeam, nil, c)
}
