package model

type Flag struct { // flag存在redis中
	MatchID     int64  // 比赛id
	TeamID      int64  // 团队id
	ContainerID string // 容器id
	Round       uint   // 比赛轮次
	Value       string // flag的值
}
