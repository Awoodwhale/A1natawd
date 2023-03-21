package v1

import (
	"github.com/gin-gonic/gin"
	"go_awd/service/team"
	"net/http"
)

// CreateTeam
// @Description: 创建team
// @param c *gin.Context
func CreateTeam(c *gin.Context) {
	var service team.CreateTeamService
	if err := c.ShouldBind(&service); err == nil {
		res := service.CreateTeam(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// UpdateTeam
// @Description: 更新team信息，只有队长可以修改信息
// @param c *gin.Context
func UpdateTeam(c *gin.Context) {
	var service team.UpdateTeamService
	if err := c.ShouldBind(&service); err == nil {
		res := service.UpdateTeam(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// LeaveTeam
// @Description: 离开团队，队长无法离开，必须先转让队长或者直接解散战队
// @param c *gin.Context
func LeaveTeam(c *gin.Context) {
	var service team.EmptyService
	if err := c.ShouldBind(&service); err == nil {
		res := service.LeaveTeam(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// DismissTeam
// @Description: 解散团队，只有队长可以解散
// @param c *gin.Context
func DismissTeam(c *gin.Context) {
	var service team.EmptyService
	if err := c.ShouldBind(&service); err == nil {
		res := service.DismissTeam(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// TransferTeam
// @Description: 转让队长，只有队长可以转让
// @param c *gin.Context
func TransferTeam(c *gin.Context) {
	var service team.TransferService
	if err := c.ShouldBind(&service); err == nil {
		res := service.TransferTeam(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// ApplyTeam
// @Description: 申请入队
// @param c *gin.Context
func ApplyTeam(c *gin.Context) {
	var service team.ApplyService
	if err := c.ShouldBind(&service); err == nil {
		res := service.ApplyTeam(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// CancelApplyTeam
// @Description: 取消申请入队
// @param c *gin.Context
func CancelApplyTeam(c *gin.Context) {
	var service team.EmptyService
	if err := c.ShouldBind(&service); err == nil {
		res := service.CancelApplyTeam(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// AcceptTeam
// @Description: 同意申请入队
// @param c *gin.Context
func AcceptTeam(c *gin.Context) {
	var service team.AcceptOrRejectService
	if err := c.ShouldBind(&service); err == nil {
		res := service.AcceptTeam(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}

// RejectTeam
// @Description: 拒绝申请入队
// @param c *gin.Context
func RejectTeam(c *gin.Context) {
	var service team.AcceptOrRejectService
	if err := c.ShouldBind(&service); err == nil {
		res := service.RejectTeam(c)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err, &service, c))
	}
}
