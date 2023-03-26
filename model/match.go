package model

type Match struct {
	Model

	GameName  string `gorm:"not null"` // 比赛名称
	Info      string `gorm:"not null"` // 比赛信息
	UserID    uint   `gorm:"not null"` // 组织比赛的用户id
	Round     uint   `gorm:"not null"` // 比赛轮次
	RoundTime uint   `gorm:"not null"` // 每轮的时间
	StartTime uint   `gorm:"not null"` // 定时开启比赛的时间
}
