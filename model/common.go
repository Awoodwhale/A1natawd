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
