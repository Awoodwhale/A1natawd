package dao

import (
	"github.com/gin-gonic/gin"
	"go_awd/model"
	"gorm.io/gorm"
)

type TokenDao struct {
	*gorm.DB
}

func NewTokenDao(ctx *gin.Context) *TokenDao {
	return &TokenDao{NewDBClient(ctx)}
}

func NewTokenDaoByDB(db *gorm.DB) *TokenDao {
	return &TokenDao{db}
}

func (dao *TokenDao) CreateRefreshToken(token *model.RefreshToken) error {
	//! 保证一个用户只有一条数据
	if _, err := dao.GetByUserID(token.UserID); err != nil {
		// 不存在user_id对应的refreshToken，新建一个
		token.ID = model.GenID() // 雪花算法生成唯一id
		return dao.DB.Create(&token).Error
	}
	// 存在就更新对应的refreshToken
	return dao.DB.Model(&model.RefreshToken{}).
		Where("user_id=?", token.UserID).
		Updates(&token).Error
}

func (dao *TokenDao) DeleteRefreshTokenByUserID(userID int64) error {
	return dao.DB.Model(&model.RefreshToken{}).
		Where("user_id=?", userID).
		Update("refresh_token", nil).
		Update("token_key", nil).Error
}

func (dao *TokenDao) GetByUserID(userID int64) (token *model.RefreshToken, err error) {
	err = dao.DB.First(&token, "user_id=?", userID).Error
	return
}

func (dao *TokenDao) GetByTokenKey(key string) (token *model.RefreshToken, err error) {
	err = dao.DB.First(&token, "token_key=?", key).Error
	return
}
