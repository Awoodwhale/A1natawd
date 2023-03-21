package model

type Challenge struct {
	Model

	Title        string  `gorm:"not null"` // 题目标题
	Info         string  // 题目描述
	DockerSource string  `gorm:"not null"`        // dockerfile存储的路径，例如/docker_source/pwn1
	BaseScore    float64 `gorm:"not null"`        // 题目分数
	Type         string  `gorm:"type:varchar(2)"` // 0表示web，1表示pwn
}
