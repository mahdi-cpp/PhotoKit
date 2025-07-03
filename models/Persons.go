package models

import "time"

type Persons struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId    int       `gorm:"references:users(id);onDelete:SET NULL" json:"userId"`
	Named     string    `json:"named"`
	CreatedAt time.Time `gorm:"default:now()" json:"createdAt"`
}
