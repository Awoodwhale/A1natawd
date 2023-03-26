package dao

import (
	"github.com/gin-gonic/gin"
	"go_awd/model"
	"gorm.io/gorm"
)

type TeamDao struct {
	*gorm.DB
}

func NewTeamDao(ctx *gin.Context) *TeamDao {
	return &TeamDao{NewDBClient(ctx)}
}

func NewTeamDaoByDB(db *gorm.DB) *TeamDao {
	return &TeamDao{db}
}

func (dao *TeamDao) CreateTeam(team *model.Team) error {
	if team.ID == 0 {
		team.ID = model.GenID()
	}
	return dao.Create(&team).Error
}

func (dao *TeamDao) GetByTeamName(teamName string) (team *model.Team, err error) {
	err = dao.First(&team, "team_name=?", teamName).Error
	return
}

func (dao *TeamDao) GetByID(id int64) (team *model.Team, err error) {
	err = dao.First(&team, "id = ?", id).Error
	return
}

func (dao *TeamDao) UpdateByID(id int64, team *model.Team) error {
	return dao.Where("id = ?", id).Save(&team).Error
}

func (dao *TeamDao) DeleteByID(id int64) error {
	return dao.Where("id = ?", id).Delete(&model.Team{}).Error
}
