package serializer

import (
	"go_awd/conf"
	"go_awd/model"
)

// User
// @Description: vo view objective	传给前端的对象
type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Avatar       string `json:"avatar"`
	Role         string `json:"role"`
	Money        uint   `json:"money"`
	Score        uint   `json:"score"`
	Sign         string `json:"sign"`
	TeamID       int64  `json:"team_id,omitempty"`
	IsTeamLeader bool   `json:"is_team_leader,omitempty"`
	CreateAt     int64  `json:"create_at"`
}

// BuildUser
// @Description: 构造一个user vo
// @param user *model.User
// @return *User
func BuildUser(user *model.User) *User {
	return &User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
		Money:        user.Money,
		Score:        user.Score,
		Sign:         user.Sign,
		TeamID:       user.TeamID,
		IsTeamLeader: user.IsTeamLeader,
		Avatar:       conf.ImgHost + ":" + conf.ImgPort + conf.ImgPath[1:] + user.Avatar,
		CreateAt:     user.CreatedAt.Unix(),
	}
}

func BuildUsers(items []*model.User) (users []*User) {
	for _, item := range items {
		users = append(users, BuildUser(item))
	}
	return
}
