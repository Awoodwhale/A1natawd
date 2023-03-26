package e

const (
	Success = 20000 // 成功
	Error   = 50000 // 失败
	Invalid = 40000 // 校验失败
)
const (
	SuccessWithRegister                 = 20001 + iota // 用户注册成功
	SuccessWithLogin                                   // 登录成功
	SuccessWithLogout                                  // 退出成功
	SuccessWithUpdateEmail                             // 修改email成功
	SuccessWithUpdatePwd                               // 修改password成功
	SuccessWithUpdate                                  // 更新成功
	SuccessWithGenCaptcha                              // 生成图形验证码成功
	SuccessWithShow                                    // 信息获取成功
	SuccessWithRecoverPwd                              // 找回密码成功
	SuccessWithCreateTeam                              // 创建团队成功
	SuccessWithLeaveTeam                               // 退出战队成功
	SuccessWithDismissTeam                             // 解散战队成功
	SuccessWithTrfTeam                                 // 转让队长成功
	SuccessWithApplyTeam                               // 申请入队成功
	SuccessWithAcceptTeam                              // 同意入队成功
	SuccessWithCancelApplyTeam                         // 取消申请入队成功
	SuccessWithRejectTeam                              // 拒绝入队成功
	SuccessWithUploadChallenge                         // 成功上传题目的build file
	SuccessWithStartTestChallenge                      // 成功开启题目测试
	SuccessWithFindStartedTestChallenge                // 成功找到已经开启的题目测试容器
	SuccessWithEndTestChallenge                        // 成功删除题目测试容器
	SuccessWithRemoveChallenge                         // 成功删除题目
)
const (
	InvalidWithExistUser                = 40001 + iota // 用户名已存在
	InvalidWithAuth                                    // 登录验证失败
	InvalidWithCaptchaKey                              // 验证码key错误
	InvalidWithCaptcha                                 // 验证码value错误
	InvalidTooManyRequest                              // 访问过多
	InvalidWithSameEmail                               // 修改邮箱时，原邮箱与目的邮箱相等
	InvalidWithPassword                                // 密码无效
	InvalidWithGenJwt                                  // 生成jwt失败
	InvalidWithShow                                    // 获取信息失败
	InvalidWithImgSize                                 // 图片过大
	InvalidWithExistTeam                               // 团队已存在
	InvalidWithCreateUser                              // 用户创建失败
	InvalidWithUpdateUser                              // 用户更新失败
	InvalidWithCreateTeam                              // 团队创建失败
	InvalidWithUpdateTeam                              // 团队更新失败
	InvalidWithFileType                                // 文件类型无效
	InvalidWithLeaderLeave                             // 团队队长不能直接退出战队
	InvalidWithLeaveTeam                               // 退出战队失败
	InvalidWithDismissTeam                             // 解散战队失败
	InvalidWithTsfTeam                                 // 转让队长失败
	InvalidWithNotExistUser                            // 用户名不存在
	InvalidWithApplyTeam                               // 申请入队失败
	InvalidWithReviewTeam                              // 团队申请审核中
	InvalidWithAcceptTeam                              // 同意入队失败
	InvalidWithCancelApplyTeam                         // 取消申请入队失败
	InvalidWithRejectTeam                              // 拒绝入队失败
	InvalidWithUploadFile                              // 文件上传失败
	InvalidWithCreateChallenge                         // 创建题目失败
	InvalidWithUpdateChallenge                         // 更新题目失败
	InvalidWithNotExistChallenge                       // 不存在的题目
	InvalidWithNotSuccessChallenge                     // 非success状态下的题目不能更新
	InvalidWithContainerInfoLost                       // 容器信息丢失
	InvalidWithNotExistStartedContainer                // 不存在已开启的容器
	InvalidWithRemoveChallenge                         // 删除题目失败
)
const (
	ErrorWithSQL        = 50001 + iota // sql错误
	ErrorWithRedis                     // redis错误
	ErrorWithEncryption                // 加密错误
	ErrorWithGenCaptcha                // 验证码生成失败
)
