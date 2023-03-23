package challenge

import (
	"go_awd/model"
)

type CreateOrUpdateChallengeImageService struct {
	Title           string `json:"title" form:"title" binding:"required,lte=50" msg:"invalid_challenge_title"`
	Info            string `json:"info" form:"info" binding:"omitempty,lte=1024" msg:"invalid_params"`
	Type            string `json:"type" form:"type" binding:"required,oneof=pwn web" msg:"invalid_challenge_type"`   // pwn | web
	BaseScore       int    `json:"base_score" form:"base_score" binding:"omitempty,gte=10" msg:"invalid_base_score"` // 基础得分
	ImageName       string `json:"image_name" form:"image_name"`                                                     // dockerhub上的镜像链接或者build image的名称
	InnerServerPort string `json:"inner_server_port" form:"inner_server_port" binding:"required,gte=0,lte=65535" msg:"invalid_port"`
}

type UpdateChallengeInfoService struct {
	Title           string `json:"title" form:"title" binding:"omitempty,lte=50"  msg:"invalid_challenge_title"`
	Info            string `json:"info" form:"info" binding:"omitempty,lte=1024" msg:"invalid_params"`
	Type            string `json:"type" form:"type" binding:"omitempty,oneof=pwn web" msg:"invalid_challenge_type"`  // pwn | web
	BaseScore       int    `json:"base_score" form:"base_score" binding:"omitempty,gte=10" msg:"invalid_base_score"` // 基础得分
	InnerServerPort string `json:"inner_server_port" form:"inner_server_port" binding:"omitempty,gte=0,lte=65535" msg:"invalid_port"`
}

type RemoveContainerService struct {
}

type EmptyService struct {
	model.BasePage
}
