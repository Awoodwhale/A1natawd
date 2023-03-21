package dao

import (
	"context"
	"go_awd/model"
	"gorm.io/gorm"
)

type UserTeamDao struct {
	*gorm.DB
}

func NewUTDao(c context.Context) *UserTeamDao {
	return &UserTeamDao{NewDBClient(c)}
}

func NewUTDaoByDB(db *gorm.DB) *UserTeamDao {
	return &UserTeamDao{db}
}

func (dao *UserTeamDao) CreateItem(ut *model.UserTeam) error {
	if ut.ID == 0 {
		ut.ID = model.GenID()
	}
	return dao.Create(&ut).Error
}

func (dao *UserTeamDao) CreateOrUpdateItemByUserID(ut *model.UserTeam) error {
	if ut.ID == 0 {
		ut.ID = model.GenID()
	}
	if ut.UserID != 0 { // 数据库可能已经存在这个item了，需要更新
		var inut *model.UserTeam
		if err := dao.First(&inut, "user_id = ?", ut.UserID).Error; err == nil {
			// 更新team
			inut.TeamID = ut.TeamID
			return dao.Save(&inut).Error
		} else {
			// 不存在这个id
			return dao.Create(&model.UserTeam{
				UserID: ut.UserID,
				TeamID: ut.TeamID,
				State:  ut.State,
			}).Error
		}
	} else { // item id为空，直接创建
		ut := &model.UserTeam{
			UserID: ut.UserID,
			TeamID: ut.TeamID,
			State:  ut.State,
		}
		return dao.Create(ut).Error
	}
}

func (dao *UserTeamDao) DeleteByUserID(userID int64) error {
	return dao.Where("user_id=?", userID).Delete(&model.UserTeam{}).Error
}

func (dao *UserTeamDao) GetUserIDsByTeamID(teamID int64) (res []int64, err error) {
	var uts []*model.UserTeam
	if err := dao.Find(&uts, "team_id = ?", teamID).Error; err != nil {
		return nil, err
	}
	for _, v := range uts {
		res = append(res, v.UserID)
	}
	return res, nil
}

func (dao *UserTeamDao) DeleteByTeamID(teamID int64) error {
	return dao.Where("team_id = ?", teamID).Delete(&model.UserTeam{}).Error
}

// ExistItemByUserID
// @Description: user是否存在ut_item
// @receiver dao *UserTeamDao
// @param id int64
// @return ut *model.UserTeam
// @return flag bool
func (dao *UserTeamDao) ExistItemByUserID(id int64) (ut *model.UserTeam, flag bool) {
	if err := dao.First(&ut, "user_id = ?", id).Error; err != nil {
		return nil, false
	}
	return ut, true
}
