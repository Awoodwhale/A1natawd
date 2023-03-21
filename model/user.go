package model

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordCost = 12
	AdminRole    = "admin"
	UserRole     = "user"
	LeaderRole   = "leader"
)

type User struct {
	Model
	Username     string `gorm:"not null;type:varchar(50)"`
	Password     string `gorm:"not null;type:varchar(100)"`
	Email        string `gorm:"not null;type:varchar(50)"`
	Role         string `gorm:"not null;type:varchar(10)"` // admin || user || none
	Avatar       string `gorm:"not null;comment:用户头像"`     // 用户头像
	Money        uint   `gorm:"not null;comment:用户金币"`
	Score        uint   `gorm:"not null;comment:用户积分"` // 用户积分
	Sign         string `gorm:"not null;type:varchar(100);comment:用户签名"`
	TeamID       int64  `gorm:"comment:用户团队id"` // 团队id
	IsTeamLeader bool   `gorm:"default:false;comment:是否是团队队长"`
}

func SetPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), PasswordCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (u *User) SetPassword(pwd string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), PasswordCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPassword(pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
	if err != nil {
		return false
	}
	return true
}
