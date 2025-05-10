package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string    `gorm:"not null;unique"`
	Name     string    `gorm:"not null"`
	BirthDay time.Time `gorm:"column:birth_day;type:date"`
	PassHash string    `gorm:"columt:pass_hash" json:"-"`
}
