package models

import "time"

type Camera struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Named     string    `json:"named"`
	UserId    int       `gorm:"references:users(id);onDelete:SET NULL" json:"userId"`
	CreatedAt time.Time `gorm:"default:now()" json:"createdAt"`
}
