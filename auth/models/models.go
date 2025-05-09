package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"not null;unique"`
	PassHash string `gorm:"pass_hash" json:"-"`
}
