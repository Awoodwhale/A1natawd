package model

type UserTeam struct {
	Model
	UserID int64 `gorm:"not null"`
	TeamID int64 `gorm:"not null"`
	State  int   `gorm:"not null;default:0"` // 0:审核中 1:成员 2:队长
}
