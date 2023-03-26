package dao

import "go_awd/model"

// migration
// @Description: 数据迁移、建表
func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			new(model.Article),
			new(model.Box),
			new(model.Challenge),
			new(model.Match),
			new(model.Team),
			new(model.User),
			new(model.RefreshToken),
			new(model.UserTeam),
			//&model.Flag{},	// 每轮的flag存在redis中
		)
	if err != nil {
		panic(err)
	}
	return
}
