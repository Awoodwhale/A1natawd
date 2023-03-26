package cache

import (
	"time"
)

// CaptchaMemStore
// @Description: 图形验证码的存储模式，存到redis中
type CaptchaMemStore struct {
	expiration time.Duration
}

func (c *CaptchaMemStore) Set(id, value string) error {
	return RedisClient.Set(DigitalCaptchaKey(id), value, c.expiration).Err()
}

func (c *CaptchaMemStore) Get(id string, clear bool) string {
	result, err := RedisClient.Get(DigitalCaptchaKey(id)).Result()
	if err != nil {
		return ""
	}
	if clear {
		RedisClient.Del(DigitalCaptchaKey(id))
	}
	return result
}

func (c *CaptchaMemStore) Verify(id, ans string, clear bool) bool {
	return c.Get(id, clear) == ans
}

func NewCaptchaMemStore(exp time.Duration) *CaptchaMemStore {
	return &CaptchaMemStore{expiration: exp}
}
