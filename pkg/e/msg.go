package e

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go_awd/pkg/util"
	"reflect"
)

var MessageFlags = map[int]string{
	Success: "success",
	Invalid: "invalid_params",
	Error:   "failed",

	SuccessWithRegister:                 "success_register",
	SuccessWithLogin:                    "success_login",
	SuccessWithLogout:                   "success_logout",
	SuccessWithGenCaptcha:               "success_gen_captcha",
	SuccessWithUpdateEmail:              "success_update_email",
	SuccessWithUpdatePwd:                "success_update_pwd",
	SuccessWithUpdate:                   "success_update",
	SuccessWithShow:                     "success_show",
	SuccessWithRecoverPwd:               "success_recover_pwd",
	SuccessWithCreateTeam:               "success_create_team",
	SuccessWithLeaveTeam:                "success_leave_team",
	SuccessWithDismissTeam:              "success_dismiss_team",
	SuccessWithTrfTeam:                  "success_transfer_team",
	SuccessWithApplyTeam:                "success_apply_team",
	SuccessWithAcceptTeam:               "success_accept_team",
	SuccessWithCancelApplyTeam:          "success_cancel_apply_team",
	SuccessWithRejectTeam:               "success_reject_team",
	SuccessWithUploadChallenge:          "success_upload_challenge",
	SuccessWithStartTestChallenge:       "success_start_test_challenge",
	SuccessWithFindStartedTestChallenge: "success_find_started_test_challenge",
	SuccessWithEndTestChallenge:         "success_end_test_challenge",
	SuccessWithRemoveChallenge:          "success_remove_challenge",
	SuccessWithBanUser:                  "success_ban_user",

	InvalidTooManyRequest:               "invalid_too_many_request",
	InvalidWithImgSize:                  "invalid_img_size",
	InvalidWithAuth:                     "invalid_auth",
	InvalidWithShow:                     "invalid_show",
	InvalidWithGenJwt:                   "invalid_gen_jwt",
	InvalidWithExistUser:                "invalid_exist_user",
	InvalidWithCaptcha:                  "invalid_captcha",
	InvalidWithPassword:                 "invalid_password",
	InvalidWithCaptchaKey:               "invalid_captcha_key",
	InvalidWithSameEmail:                "invalid_same_email",
	InvalidWithExistTeam:                "invalid_exist_team",
	InvalidWithCreateUser:               "invalid_create_user",
	InvalidWithUpdateUser:               "invalid_update_user",
	InvalidWithCreateTeam:               "invalid_create_team",
	InvalidWithUpdateTeam:               "invalid_update_team",
	InvalidWithFileType:                 "invalid_file_type",
	InvalidWithLeaderLeave:              "invalid_leader_leave",
	InvalidWithLeaveTeam:                "invalid_leave_team",
	InvalidWithDismissTeam:              "invalid_dismiss_team",
	InvalidWithTsfTeam:                  "invalid_transfer_team",
	InvalidWithNotExistUser:             "invalid_not_exist_user",
	InvalidWithApplyTeam:                "invalid_apply_team",
	InvalidWithReviewTeam:               "invalid_review_team",
	InvalidWithAcceptTeam:               "invalid_accept_team",
	InvalidWithCancelApplyTeam:          "invalid_cancel_apply_team",
	InvalidWithRejectTeam:               "invalid_reject_team",
	InvalidWithUploadFile:               "invalid_upload_file",
	InvalidWithCreateChallenge:          "invalid_create_challenge",
	InvalidWithUpdateChallenge:          "invalid_update_challenge",
	InvalidWithNotExistChallenge:        "invalid_not_exist_challenge",
	InvalidWithNotSuccessChallenge:      "invalid_not_success_challenge",
	InvalidWithContainerInfoLost:        "invalid_container_info_lost",
	InvalidWithNotExistStartedContainer: "invalid_not_exist_started_container",
	InvalidWithRemoveChallenge:          "invalid_remove_challenge",
	InvalidWithBanUser:                  "invalid_ban_user",

	ErrorWithSQL:        "error_sql",
	ErrorWithRedis:      "error_redis",
	ErrorWithEncryption: "error_encryption",
	ErrorWithGenCaptcha: "error_gen_captcha",
}

// GetMessageByCode
// @Description: 获取code对应的message
// @param code int
// @return string
func GetMessageByCode(code int, c *gin.Context) string {
	if key, ok := MessageFlags[code]; ok {
		return util.GetI18nMsg(key, c)
	}
	return util.GetI18nMsg("failed", c)
}

// HandleBindingError
// @Description: 获取binding错误的msg tag
// @param err error
// @param obj any
// @param c *gin.Context
// @return string
func HandleBindingError(err error, obj any, c *gin.Context) string {
	if err == nil {
		return ""
	}
	if errs, ok := err.(validator.ValidationErrors); ok {
		getObj := reflect.TypeOf(obj)
		for _, v := range errs {
			if f, exist := getObj.Elem().FieldByName(v.Field()); exist {
				return util.GetI18nMsg(f.Tag.Get("msg"), c)
			}
		}
	}
	return err.Error()
}
