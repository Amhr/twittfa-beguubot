package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func NewGorm() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(mariadb:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MARIADB_USER"), os.Getenv("MARIADB_PASSWD"), os.Getenv("MARIADB_DATABASE"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database %w", err)
	}
	MigrateGorm(db)
	return db, nil
}
