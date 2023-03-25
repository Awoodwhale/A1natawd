package model

// Challenge
// @Description: 题目model
type Challenge struct {
	Model

	Title           string `gorm:"not null;unique"`           // 题目标题
	Info            string `gorm:"not null"`                  // 题目描述
	Type            string `gorm:"not null;type:varchar(10)"` // web | pwn
	BaseScore       int    `gorm:"not null"`                  // 题目分数
	ImageName       string `gorm:"not null"`                  // 题目的docker image name
	InnerServerPort string `gorm:"not null"`                  // 容器内部题目开放的端口
	State           string `gorm:"not null;default:error"`    // 是否build image成功, error | success | building
}
