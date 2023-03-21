package serializer

import (
	"go_awd/conf"
	"go_awd/model"
)

type Team struct {
	ID       int64  `json:"id"`
	LeaderID int64  `json:"leader_id"`
	TeamName string `json:"team_name"`
	Avatar   string `json:"avatar"`
	Score    uint   `json:"score"`
	Rank     uint   `json:"rank"`
	Sign     string `json:"sign"`
	CreateAt int64  `json:"create_at"`
}

func BuildTeam(team *model.Team) *Team {
	return &Team{
		ID:       team.ID,
		TeamName: team.TeamName,
		Avatar:   conf.ImgHost + ":" + conf.ImgPort + conf.ImgPath[1:] + team.Avatar,
		Sign:     team.Sign,
		Score:    team.Score,
		Rank:     team.Rank,
		LeaderID: team.LeaderID,
		CreateAt: team.CreatedAt.Unix(),
	}
}
