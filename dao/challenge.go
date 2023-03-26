package dao

import (
	"github.com/gin-gonic/gin"
	"go_awd/model"
	"gorm.io/gorm"
)

type ChallengeDao struct {
	*gorm.DB
}

func NewChallengeDao(ctx *gin.Context) *ChallengeDao {
	return &ChallengeDao{NewDBClient(ctx)}
}

func NewChallengeDaoByDB(db *gorm.DB) *ChallengeDao {
	return &ChallengeDao{db}
}

func (dao *ChallengeDao) CreateOrUpdateChallenge(chal *model.Challenge) error {
	if chal.ID == 0 {
		chal.ID = model.GenID()
	}
	if ch, err := dao.GetByTitle(chal.Title); err == nil {
		// 存在就更新
		chal.ID = ch.ID
		return dao.UpdateByTitle(chal.Title, chal)
	}
	// 不存在就创建
	return dao.DB.Create(&chal).Error
}

func (dao *ChallengeDao) GetByTitle(title string) (chal *model.Challenge, err error) {
	err = dao.DB.First(&chal, "title = ?", title).Error
	return
}

func (dao *ChallengeDao) UpdateByTitle(title string, chal *model.Challenge) error {
	return dao.DB.Where("title = ?", title).Updates(&chal).Error
}

func (dao *ChallengeDao) ListByCondition(condition map[string]any, page *model.BasePage) (chals []*model.Challenge, err error) {
	err = dao.DB.Where(condition).
		Offset(int((page.PageNum - 1) * page.PageSize)). // 分页
		Limit(int(page.PageSize)).Find(&chals).Error
	return
}

func (dao *ChallengeDao) GetByID(id int64) (chal *model.Challenge, err error) {
	err = dao.DB.First(&chal, "id = ?", id).Error
	return
}

func (dao *ChallengeDao) UpdateByID(chal *model.Challenge) error {
	return dao.DB.Where("id = ?", chal.ID).Updates(&chal).Error
}

func (dao *ChallengeDao) DeleteByID(id int64) error {
	return dao.DB.Delete(&model.Challenge{}, "id = ?", id).Error
}
