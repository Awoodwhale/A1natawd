package util

import (
	ginI18n "github.com/fishjar/gin-i18n"
	"github.com/gin-gonic/gin"
)

func GetI18nMsg(key string, c *gin.Context) string {
	return c.MustGet("Localizer").(*ginI18n.UserLocalize).GetMsg(key)
}
