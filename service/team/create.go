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

// CreateTeam
// @Description: 创建团队
// @receiver t *CreateTeamService
// @param c *gin.Context
// @return serializer.Response
func (t *CreateTeamService) CreateTeam(c *gin.Context) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	userDao := dao.NewUserDao(c)
	user, err := userDao.GetUserByID(claims.ID)
	if err != nil {
		return serializer.RespCode(e.InvalidWithAuth, c)
	}
	if user.TeamID != 0 { // 已存在团队，无法创建
		return serializer.RespCode(e.InvalidWithExistTeam, c)
	}

	teamDao := dao.NewTeamDaoByDB(userDao.DB)
	if _, err := teamDao.GetByTeamName(t.TeamName); err == nil {
		// 团队名重复
		return serializer.RespCode(e.InvalidWithExistTeam, c)
	}

	teamID := model.GenID() // 生成id
	teamAvatar := "avatar.png"
	// 上传团队头像
	file, err := c.FormFile("file")
	if err == nil {
		if file.Size > conf.ImgMaxSize { // 图片过大
			return serializer.RespCode(e.InvalidWithImgSize, c)
		}
		// 上传图片到本地
		imgPath := conf.ImgPath + "team/" + strconv.FormatInt(teamID, 10) + "/" + file.Filename
		if err := c.SaveUploadedFile(file, imgPath); err != nil {
			service.Debugln("upload image file error,", err.Error())
			return serializer.RespCode(e.InvalidWithUpdateUser, c)
		}
		teamAvatar = "team/" + strconv.FormatInt(teamID, 10) + "/" + file.Filename
	}
	// 创建team到MySQL
	team := &model.Team{
		LeaderID: user.ID,
		Avatar:   teamAvatar,
		TeamName: t.TeamName,
		Sign:     t.Sign,
		Score:    0,
		Rank:     0,
	}
	team.ID = teamID
	if err := teamDao.CreateTeam(team); err != nil {
		service.Errorln("create team mysql error,", err.Error())
		return serializer.RespCode(e.InvalidWithCreateTeam, c)
	}

	// 创建用户-团队表的item
	utDao := dao.NewUTDaoByDB(userDao.DB)
	if err := utDao.CreateItem(&model.UserTeam{
		UserID: user.ID,
		TeamID: teamID,
		State:  2, // 职责是队长
	}); err != nil {
		service.Errorln("create user_team mysql error,", err.Error())
		return serializer.RespCode(e.InvalidWithCreateTeam, c)
	}

	// 更新用户的队长状态
	user.TeamID = teamID
	user.IsTeamLeader = true
	if err := userDao.UpdateByID(user.ID, user); err != nil {
		service.Errorln("edit user isLeader state mysql error,", err.Error())
		return serializer.RespCode(e.InvalidWithCreateTeam, c)
	}
	return serializer.RespSuccess(e.SuccessWithCreateTeam, serializer.BuildTeam(team), c)
}
