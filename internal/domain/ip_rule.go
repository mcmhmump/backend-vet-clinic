package domain

import "time"

type IPRule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ListType  string    `gorm:"size:20;not null" json:"list_type"`
	Value     string    `gorm:"size:100;not null" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
