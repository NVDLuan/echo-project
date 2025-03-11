package auth

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	Username     string     `json:"username"`
	Email        string     `json:"email" gorm:"unique"`
	PasswordHash string     `json:"-"`
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
