package model

type Flag struct { // flag存在redis中
	TeamID  uint   // 团队id
	MatchID uint   // 比赛id
	BoxID   uint   // 容器id
	Round   uint   // 比赛轮次
	Value   string // flag的值
}
