package model

import (
	"github.com/bwmarrin/snowflake"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        int64 `gorm:"column:id;type:bigint;primary_key;auto_increment:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

var node, _ = snowflake.NewNode(1)

func GenID() int64 {
	return node.Generate().Int64()
}

// BasePage
// @Description: 分页model
type BasePage struct {
	PageNum  uint `json:"page_num" form:"page_num"`
	PageSize uint `json:"page_size" form:"page_size"`
}
