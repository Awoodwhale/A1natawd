package dao

import (
	"github.com/gin-gonic/gin"
	"go_awd/model"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

// NewUserDao
// @Description: 通过ctx获取userDao
// @param ctx context.Context
// @return *UserDao
func NewUserDao(ctx *gin.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

// NewUserDaoByDB
// @Description: 通过db获取userDao
// @param db *gorm.DB
// @return *UserDao
func NewUserDaoByDB(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

// ExistOrNotByUserNameOrEmail
// @Description: 通过username|| email查询用户是否存在
// @receiver dao *UserDao
// @param name string
// @param email string
// @return user *model.User
// @return exist bool
func (dao *UserDao) ExistOrNotByUserNameOrEmail(name, email string) (user *model.User, exist bool) {
	err := dao.DB.First(&user, "username = ? OR email = ?", name, email).Error
	if user == nil || err == gorm.ErrRecordNotFound { // 不存在
		return nil, false
	} else if user != nil {
		return user, true
	}
	return nil, false // 不存在
}

// CreateUser
// @Description: 创建用户
// @receiver dao *UserDao
// @param user *model.User
// @return error
func (dao *UserDao) CreateUser(user *model.User) error {
	if user.ID == 0 {
		user.ID = model.GenID() // 雪花算法生成唯一id
	}
	return dao.DB.Create(&user).Error
}

// GetUserByUsername
// @Description: 通过用户名获取user
// @receiver dao *UserDao
// @param username string
// @return user *model.User
// @return err error
func (dao *UserDao) GetUserByUsername(username string) (user *model.User, err error) {
	err = dao.DB.First(&user, "username = ?", username).Error
	return
}

// GetUserByEmail
// @Description: 通过email获取用户
// @receiver dao *UserDao
// @param email string
// @return user *model.User
// @return err error
func (dao *UserDao) GetUserByEmail(email string) (user *model.User, err error) {
	err = dao.DB.First(&user, "email = ?", email).Error
	return
}

func (dao *UserDao) GetUserByID(id int64) (user *model.User, err error) {
	err = dao.DB.First(&user, "id = ?", id).Error
	return
}

// UpdatePwdByUsername
// @Description: 通过用户名修改密码
// @receiver dao *UserDao
// @param username string
// @param password string
// @return error
func (dao *UserDao) UpdatePwdByUsername(username string, password string) error {
	return dao.DB.Model(&model.User{}).Where("username = ?", username).Update("password", password).Error
}

// UpdateEmailByUsername
// @Description: 通过用户名修改email
// @receiver dao *UserDao
// @param username string
// @param email string
// @return error
func (dao *UserDao) UpdateEmailByUsername(username string, email string) error {
	return dao.DB.Model(&model.User{}).
		Where("username = ?", username).
		Update("email", email).Error
}

func (dao *UserDao) UpdateByID(id int64, user *model.User) error {
	return dao.DB.Model(&model.User{}).
		Where("id = ?", id).
		Save(&user).Error
}

func (dao *UserDao) GetUserByUsernameAndEmail(username string, email string) (user *model.User, err error) {
	err = dao.DB.First(&user, "username = ? AND email = ?", username, email).Error
	return
}

func (dao *UserDao) ListByCondition(condition map[string]any, page *model.BasePage) (users []*model.User, err error) {
	err = dao.DB.Where(condition).
		Offset(int((page.PageNum - 1) * page.PageSize)).
		Limit(int(page.PageSize)).Find(&users).Error
	return
}
