package domain

import "time"

type AccessLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ClientIP  string    `gorm:"size:100;not null" json:"client_ip"`
	URL       string    `gorm:"size:255;not null" json:"url"`
	Allowed   bool      `json:"allowed"`
	Reason    string    `gorm:"size:100" json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}
