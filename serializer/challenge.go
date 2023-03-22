package serializer

import "go_awd/model"

type Challenge struct {
	ID              int64  `json:"id"`
	Title           string `json:"title"`             // 题目标题
	Info            string `json:"info"`              // 题目描述
	Type            string `json:"type"`              // web | pwn
	BaseScore       int    `json:"base_score"`        // 题目分数
	ImageName       string `json:"image_name"`        // 镜像名称
	InnerServerPort string `json:"inner_server_port"` // 容器内部题目开放的端口
	State           string `json:"state"`             // 题目状态 building | error | success
	CreateAt        int64  `json:"create_at"`
}

func BuildChallenge(chal *model.Challenge) *Challenge {
	return &Challenge{
		ID:              chal.ID,
		Title:           chal.Title,
		Info:            chal.Info,
		Type:            chal.Type,
		BaseScore:       chal.BaseScore,
		ImageName:       chal.ImageName,
		InnerServerPort: chal.InnerServerPort,
		State:           chal.State,
		CreateAt:        chal.CreatedAt.Unix(),
	}
}
func BuildChallenges(items []*model.Challenge) (chals []*Challenge) {
	for _, item := range items {
		chals = append(chals, BuildChallenge(item))
	}
	return
}
