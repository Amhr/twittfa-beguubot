package models

import "gorm.io/gorm"

func MigrateGorm(g *gorm.DB) {
	g.AutoMigrate(&DBUser{})
	g.AutoMigrate(&DBAnnmsg{})
}
