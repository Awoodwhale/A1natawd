package v1

import (
	"github.com/gin-gonic/gin"
	"go_awd/pkg/e"
	"go_awd/serializer"
	"go_awd/service/challenge"
	"net/http"
	"strconv"
)

// CreateOrUpdateChallenge
// @Description: 管理员创建题目
// @param c *gin.Context
func CreateOrUpdateChallenge(c *gin.Context) {
	var service challenge.CreateOrUpdateChallengeImageService
	if err := c.ShouldBind(&service); err == nil {
		res := service.CreateOrUpdateChallenge(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// ShowChallenges
// @Description: 获取题目列表
// @param c *gin.Context
func ShowChallenges(c *gin.Context) {
	var service challenge.EmptyService
	if err := c.ShouldBind(&service); err == nil {
		res := service.ShowChallenges(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// UpdateChallengeInfo
// @Description: 更新题目描述信息
// @param c *gin.Context
func UpdateChallengeInfo(c *gin.Context) {
	var service challenge.UpdateChallengeInfoService
	if err := c.ShouldBind(&service); err == nil {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, serializer.RespCode(e.Invalid, c))
			return
		}
		res := service.UpdateChallenge(c, id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// StartTestChallenge
// @Description: 开启题目测试容器
// @param c *gin.Context
func StartTestChallenge(c *gin.Context) {
	var service challenge.EmptyService
	if err := c.ShouldBind(&service); err == nil {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, serializer.RespCode(e.Invalid, c))
			return
		}
		res := service.StartTestChallenge(c, id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// EndTestChallenge
// @Description: 关闭题目测试容器
// @param c *gin.Context
func EndTestChallenge(c *gin.Context) {
	var service challenge.EmptyService
	if err := c.ShouldBind(&service); err == nil {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, serializer.RespCode(e.Invalid, c))
			return
		}
		res := service.EndTestChallenge(c, id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// RemoveChallenge
// @Description: 删除题目
// @param c *gin.Context
func RemoveChallenge(c *gin.Context) {
	var service challenge.EmptyService
	if err := c.ShouldBind(&service); err == nil {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, serializer.RespCode(e.Invalid, c))
			return
		}
		res := service.RemoveChallenge(c, id)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}
