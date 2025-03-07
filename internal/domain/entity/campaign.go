package entity

import (
	"time"
)

// Campaign はキャンペーン情報を表すエンティティです
type Campaign struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	AccountID uint      `json:"account_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Budget    float64   `json:"budget"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
