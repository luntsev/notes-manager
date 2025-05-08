package database

import "gorm.io/gorm"

type PostgresDB struct {
	Db *gorm.DB
}
