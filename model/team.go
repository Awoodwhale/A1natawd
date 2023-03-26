package model

type Team struct {
	Model
	TeamName string `gorm:"not null;type:varchar(50)"`
	Avatar   string `gorm:"not null;type:varchar(100)"` // 团队头像
	Sign     string `gorm:"not null;type:varchar(100)"` // 团队签名
	Score    uint   // 团队积分
	Rank     uint   // 团队排名
	LeaderID int64  `gorm:"not null"`
}
