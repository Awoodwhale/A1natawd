package model

type Box struct {
	Model

	TeamID      uint   `gorm:"not null"` // 团队id
	MatchID     uint   `gorm:"not null"` // 比赛id
	IP          string `gorm:"not null"` // 容器的ip（假设可以部署多个服务器，那么可以分配ip）
	Port        uint   `gorm:"not null"` // 容器的题目端口
	SSHPort     uint   `gorm:"not null"` // ssh的连接端口
	SSHUserName string `gorm:"not null"` // ssh的用户名
	SSHPassword string `gorm:"not null"` // ssh的密码
	IsDown      bool   `gorm:"not null"` // 容器是否宕机
	Info        string // 容器描述
}
