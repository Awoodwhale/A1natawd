package challenge

import (
	"github.com/gin-gonic/gin"
	"go_awd/cache"
	"go_awd/dao"
	"go_awd/pkg/e"
	"go_awd/pkg/wdocker"
	wjwt "go_awd/pkg/wjwt"
	"go_awd/serializer"
	"go_awd/service"
)

// EndTestChallenge
// @Description: 结束challenge的环境
// @receiver s *EmptyService
// @param c *gin.Context
// @param id int64
// @return serializer.Response
func (s *EmptyService) EndTestChallenge(c *gin.Context, id int64) serializer.Response {
	claims := c.MustGet("claims").(*wjwt.Claims)
	exist, err := cache.RedisClient.SIsMember(cache.AdminStartTestChallengeKey(claims.ID), id).Result()
	if err != nil {
		return serializer.RespCode(e.InvalidWithContainerInfoLost, c)
	}
	if !exist {
		return serializer.RespCode(e.InvalidWithNotExistStartedContainer, c)
	}
	containerID, err := cache.RedisClient.HGet(cache.AdminContainerKey(claims.ID, id), "container_id").Result()
	if err != nil {
		return serializer.RespCode(e.InvalidWithContainerInfoLost, c)
	}
	go func() {
		if err := wdocker.NewDockerClient().RemoveContainer(containerID); err != nil { // 删除容器
			service.Errorln("EndTestChallenge RemoveContainer error,", err.Error())
		}
	}()
	// 删除redis中存储的信息
	cache.RedisClient.Del(cache.AdminContainerKey(claims.ID, id))
	cache.RedisClient.SRem(cache.AdminStartTestChallengeKey(claims.ID), id)
	return serializer.RespSuccess(e.SuccessWithEndTestChallenge, nil, c)
}

func (s *EmptyService) RemoveChallenge(c *gin.Context, id int64) serializer.Response {
	chalDao := dao.NewChallengeDao(c)
	chal, err := chalDao.GetByID(id)
	if err != nil {
		return serializer.RespCode(e.InvalidWithNotExistChallenge, c)
	}
	cli := wdocker.NewDockerClient()
	if err := cli.RemoveImage(chal.ImageName); err != nil {
		return serializer.RespCode(e.InvalidWithRemoveChallenge, c)
	}

	if err := chalDao.DeleteByID(chal.ID); err != nil {
		return serializer.RespCode(e.InvalidWithRemoveChallenge, c)
	}

	return serializer.RespSuccess(e.SuccessWithRemoveChallenge, nil, c)
}
