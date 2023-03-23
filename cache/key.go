package cache

import "fmt"

func ClientOperationIPKey(op, ip string) string {
	return fmt.Sprintf("client:ip:%s:%s", op, ip)
}

func EmailCaptchaKey(typ, email string) string {
	return fmt.Sprintf("user:captcha:email:%s:%s", typ, email)
}

func EmailSendCountKey(typ, email string) string {
	return fmt.Sprintf("user:limit:email:%s:%s", typ, email)
}

func DigitalCaptchaKey(id string) string {
	return fmt.Sprintf("user:captcha:digital:%s", id)
}

func StoreAccessTokenKeyKey(userID int64) string {
	// 获取userID对应的用户存储的accessToken的tokenKey
	return fmt.Sprintf("user:token:store:%v", userID)
}

func AccessTokenKey(key string) string {
	// 获取tokenKey所存储的accessToken
	return fmt.Sprintf("user:token:access:%s", key)
}

// AdminStartTestChallengeKey
// @Description: 返回管理员已经开启的题目的key，在redis中存储的是一个list，存放所有开启的题目
// @param adminID int64
// @return string
func AdminStartTestChallengeKey(adminID int64) string {
	return fmt.Sprintf("admin:%v:challenge", adminID)
}

// AdminContainerKey
// @Description: 返回管理员开启的containerID，一个challengeID对应一个containerID
// @param challengeID int64
// @return string
func AdminContainerKey(adminID, challengeID int64) string {
	return fmt.Sprintf("admin:%v:container:%v", adminID, challengeID)
}
