package model

type RefreshToken struct {
	Model
	UserID       int64 `gorm:"not null"`
	TokenKey     string
	RefreshToken string
}
