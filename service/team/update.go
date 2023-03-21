package team

import (
	"github.com/gin-gonic/gin"
	"go_awd/conf"
	"go_awd/dao"
	"go_awd/model"
	"go_awd/pkg/e"
	wjwt "go_awd/pkg/wjwt"
	"go_awd/serializer"
	"go_awd/service"
	"strconv"
)

// UpdateTeam
// @Description: 修改团队信息，只有队长可以
// @receiver t *UpdateTeamService
// @param c *gin.Context
// @return serializer.Response
func (t *UpdateTeamService) UpdateTeam(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(claims.ID)
	if err != nil || user.TeamID == 0 || !user.IsTeamLeader { // 不存在user或者user没有team不能修改
		return serializer.RespCode(e.InvalidWithAuth, c)
	}
	teamDao := dao.NewTeamDaoByDB(userDao.DB)
	if tt, err := teamDao.GetByTeamName(t.TeamName); err == nil {
		// 不是同一个领队，但是想要修改名称为别人的团队名
		if tt.LeaderID != user.ID {
			return serializer.RespCode(e.InvalidWithExistTeam, c)
		}
	}
	team, err := teamDao.GetByID(user.TeamID)
	if err != nil || team.LeaderID != user.ID { // team获取失败或者该用户不是队长那么无法修改
		return serializer.RespCode(e.InvalidWithUpdateTeam, c)
	}
	// 处理update逻辑
	file, err := c.FormFile("file")
	if err == nil {
		if file.Size > conf.ImgMaxSize { // 图片过大
			return serializer.RespCode(e.InvalidWithImgSize, c)
		}
		// 上传图片到本地
		imgPath := conf.ImgPath + "team/" + strconv.FormatInt(team.ID, 10) + "/" + file.Filename
		if err := c.SaveUploadedFile(file, imgPath); err != nil {
			service.Debugln("upload image file error,", err.Error())
			return serializer.RespCode(e.InvalidWithUpdateUser, c)
		}
		team.Avatar = "team/" + strconv.FormatInt(team.ID, 10) + "/" + file.Filename
	}

	if t.TeamName != "" {
		team.TeamName = t.TeamName
	}
	if t.Sign != "" {
		team.Sign = t.Sign
	}

	if err := teamDao.UpdateByID(team.ID, team); err != nil {
		service.Errorln("update team to mysql error,", err.Error())
		return serializer.RespCode(e.InvalidWithUpdateTeam, c)
	}

	return serializer.RespSuccess(e.SuccessWithUpdate, serializer.BuildTeam(team), c)
}

// LeaveTeam
// @Description: 队员离开战队
// @receiver t *EmptyService
// @param c *gin.Context
// @return serializer.Response
func (t *EmptyService) LeaveTeam(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(claims.ID)
	if err != nil || user.TeamID == 0 {
		// 不存在team无法退出
		return serializer.RespCode(e.InvalidWithAuth, c)
	}
	// 退出团队，如果是队长，不能直接退出，必须先转让团队或者解散战队
	if user.IsTeamLeader {
		return serializer.RespCode(e.InvalidWithLeaderLeave, c)
	}
	// 是普通成员那么直接退队
	utDao := dao.NewUTDaoByDB(userDao.DB)
	if err := utDao.DeleteByUserID(user.ID); err != nil { // 删除ut表中的记录
		return serializer.RespCode(e.InvalidWithLeaveTeam, c)
	}
	user.TeamID = 0
	user.IsTeamLeader = false
	if err := userDao.UpdateByID(user.ID, user); err != nil { // 修改user信息，置空team_id
		return serializer.RespCode(e.InvalidWithLeaveTeam, c)
	}

	return serializer.RespSuccess(e.SuccessWithLeaveTeam, nil, c)
}

// DismissTeam
// @Description: 解散团队
// @receiver t *EmptyService
// @param c *gin.Context
// @return serializer.Response
func (t *EmptyService) DismissTeam(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(claims.ID)
	if err != nil || user.TeamID == 0 || !user.IsTeamLeader {
		// 不存在team或者不是队长无法解散
		return serializer.RespCode(e.InvalidWithAuth, c)
	}
	// 先去找到所有的战队成员，把他们剔除战队
	utDao := dao.NewUTDaoByDB(userDao.DB)
	userIDs, err := utDao.GetUserIDsByTeamID(user.TeamID)
	if err != nil {
		service.Errorln("dismiss team mysql error, ", err.Error())
		return serializer.RespCode(e.InvalidWithDismissTeam, c)
	}
	for _, v := range userIDs {
		// 每一个用户都退出
		iu, err := userDao.GetUserByID(v)
		if err != nil {
			service.Errorln("dismiss team user error,", err.Error())
			continue // 某个队员退出失败
		}
		iu.TeamID = 0
		iu.IsTeamLeader = false
		if err := userDao.UpdateByID(iu.ID, iu); err != nil {
			service.Errorln("dismiss team user error,", err.Error())
			continue // 退出失败
		}
	}
	// 删除ut的数据
	if err := utDao.DeleteByTeamID(user.TeamID); err != nil {
		return serializer.RespCode(e.InvalidWithDismissTeam, c)
	}
	// 删除团队数据
	teamDao := dao.NewTeamDaoByDB(userDao.DB)
	if err := teamDao.DeleteByID(user.TeamID); err != nil {
		return serializer.RespCode(e.InvalidWithDismissTeam, c)
	}
	return serializer.RespSuccess(e.SuccessWithDismissTeam, nil, c)
}

// TransferTeam
// @Description: 转让队长
// @receiver t *TransferService
// @param c *gin.Context
// @return serializer.Response
func (t *TransferService) TransferTeam(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(claims.ID)
	if err != nil || user.TeamID == 0 || !user.IsTeamLeader || user.ID == t.UserID {
		// 不存在team或者不是队长或者想自己转给自己，都无法转让队长
		return serializer.RespCode(e.InvalidWithAuth, c)
	}
	// 获取新队长的信息
	newLeader, err := userDao.GetUserByID(t.UserID)
	if err != nil {
		return serializer.RespCode(e.InvalidWithNotExistUser, c)
	}
	newLeader.TeamID = user.TeamID
	newLeader.IsTeamLeader = true

	// 把队长任命给别人，更新team的leader_id
	teamDao := dao.NewTeamDaoByDB(userDao.DB)
	if err := teamDao.UpdateByID(user.TeamID, &model.Team{LeaderID: t.UserID}); err != nil {
		service.Errorln("更新team表的leader_id error,", err.Error())
		return serializer.RespCode(e.InvalidWithTsfTeam, c)

	}

	// ut表进行删除、更新
	utDao := dao.NewUTDaoByDB(userDao.DB)
	if err := utDao.DeleteByUserID(user.ID); err != nil {
		service.Errorln("ut表进行删除 error,", err.Error())
		return serializer.RespCode(e.InvalidWithTsfTeam, c)
	}
	if err := utDao.CreateOrUpdateItemByUserID(&model.UserTeam{
		UserID: t.UserID,
		TeamID: user.TeamID,
		State:  2,
	}); err != nil {
		service.Errorln("ut表进行更新 error,", err.Error())
		return serializer.RespCode(e.InvalidWithTsfTeam, c)
	}
	user.TeamID = 0
	user.IsTeamLeader = false

	// 更新前队长的信息
	if err := userDao.UpdateByID(user.ID, user); err != nil {
		service.Errorln("user表进行前队长更新 error,", err.Error())
		return serializer.RespCode(e.InvalidWithTsfTeam, c)
	}

	// 更新后队长的信息
	if err := userDao.UpdateByID(newLeader.ID, newLeader); err != nil {
		service.Errorln("user表进行新队长更新 error,", err.Error())
		return serializer.RespCode(e.InvalidWithTsfTeam, c)
	}

	return serializer.RespSuccess(e.SuccessWithTrfTeam, nil, c)
}
