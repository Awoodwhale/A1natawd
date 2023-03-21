package team

// CreateTeamService
// @Description: 创建团队service
type CreateTeamService struct {
	TeamName string `json:"team_name" form:"team_name" binding:"required,lte=50" msg:"invalid_team_name"`
	Sign     string `json:"sign" form:"sign" binding:"omitempty,lte=100" msg:"invalid_sign"`
}

// UpdateTeamService
// @Description: 更新团队service
type UpdateTeamService struct {
	TeamName string `json:"team_name" form:"team_name" binding:"omitempty,lte=50" msg:"invalid_team_name"`
	Sign     string `json:"sign" form:"sign" binding:"omitempty,lte=100" msg:"invalid_sign"`
}

// EmptyService
// @Description: empty service
type EmptyService struct{}

// TransferService
// @Description: 转让团队队长职位service
type TransferService struct {
	UserID int64 `json:"user_id" form:"user_id" binding:"required,lte=20" msg:"invalid_params"`
}

// ApplyService
// @Description: 申请加入战队的service
type ApplyService struct {
	TeamID int64 `json:"team_id" form:"team_id" binding:"required,lte=20" msg:"invalid_params"`
}

// AcceptOrRejectService
// @Description: 同意或拒绝加入战队的service
type AcceptOrRejectService struct {
	UserID int64 `json:"user_id" form:"user_id" binding:"required,lte=20" msg:"invalid_params"`
}
