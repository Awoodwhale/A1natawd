package model

type Article struct {
	Model
	Title  string `gorm:"not null"` // 公告标题
	Info   string `gorm:"not null"` // 公告内容
	UserID uint   `gorm:"not null"` // 发布公告的用户id
}
