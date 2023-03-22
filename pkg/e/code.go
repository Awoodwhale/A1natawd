package e

const (
	Success = 20000 // 成功
	Error   = 50000 // 失败
	Invalid = 40000 // 校验失败

	SuccessWithRegister           = 20001 // 用户注册成功
	SuccessWithLogin              = 20002 // 登录成功
	SuccessWithLogout             = 20003 // 退出成功
	SuccessWithUpdateEmail        = 20004 // 修改email成功
	SuccessWithUpdatePwd          = 20005 // 修改password成功
	SuccessWithUpdate             = 20006 // 更新成功
	SuccessWithGenCaptcha         = 20007 // 生成图形验证码成功
	SuccessWithShow               = 20008 // 信息获取成功
	SuccessWithRecoverPwd         = 20009 // 找回密码成功
	SuccessWithCreateTeam         = 20010 // 创建团队成功
	SuccessWithLeaveTeam          = 20011 // 退出战队成功
	SuccessWithDismissTeam        = 20012 // 解散战队成功
	SuccessWithTrfTeam            = 20013 // 转让队长成功
	SuccessWithApplyTeam          = 20014 // 申请入队成功
	SuccessWithAcceptTeam         = 20015 // 同意入队成功
	SuccessWithCancelApplyTeam    = 20016 // 取消申请入队成功
	SuccessWithRejectTeam         = 20017 // 拒绝入队成功
	SuccessWithUploadChallenge    = 20018 // 成功上传题目的build file
	SuccessWithStartTestChallenge = 20019 // 成功开启题目测试

	InvalidWithExistUser           = 40001 // 用户名已存在
	InvalidWithAuth                = 40002 // 登录验证失败
	InvalidWithCaptchaKey          = 40003 // 验证码key错误
	InvalidWithCaptcha             = 40004 // 验证码value错误
	InvalidTooManyRequest          = 40005 // 访问过多
	InvalidWithSameEmail           = 40006 // 修改邮箱时，原邮箱与目的邮箱相等
	InvalidWithPassword            = 40007 // 密码无效
	InvalidWithGenJwt              = 40008 // 生成jwt失败
	InvalidWithShow                = 40009 // 获取信息失败
	InvalidWithImgSize             = 40010 // 图片过大
	InvalidWithExistTeam           = 40011 // 团队已存在
	InvalidWithCreateUser          = 40012 // 用户创建失败
	InvalidWithUpdateUser          = 40013 // 用户更新失败
	InvalidWithCreateTeam          = 40014 // 团队创建失败
	InvalidWithUpdateTeam          = 40015 // 团队更新失败
	InvalidWithFileType            = 40016 // 文件类型无效
	InvalidWithLeaderLeave         = 40017 // 团队队长不能直接退出战队
	InvalidWithLeaveTeam           = 40018 // 退出战队失败
	InvalidWithDismissTeam         = 40019 // 解散战队失败
	InvalidWithTsfTeam             = 40020 // 转让队长失败
	InvalidWithNotExistUser        = 40021 // 用户名不存在
	InvalidWithApplyTeam           = 40022 // 申请入队失败
	InvalidWIthReviewTeam          = 40023 // 团队申请审核中
	InvalidWithAcceptTeam          = 40024 // 同意入队失败
	InvalidWithCancelApplyTeam     = 40025 // 取消申请入队失败
	InvalidWithRejectTeam          = 40026 // 拒绝入队失败
	InvalidWIthUploadFile          = 40027 // 文件上传失败
	InvalidWithCreateChallenge     = 40028 // 创建题目失败
	InvalidWIthUpdateChallenge     = 40029 // 更新题目失败
	InvalidWIthNotExistChallenge   = 40030 // 不存在的题目
	InvalidWithNotSuccessChallenge = 40031 // 非success状态下的题目不能更新

	ErrorWithSQL        = 50001 // sql错误
	ErrorWithRedis      = 50002 // redis错误
	ErrorWithEncryption = 50004 // 加密错误
	ErrorWithGenCaptcha = 50005 // 验证码生成失败
)
