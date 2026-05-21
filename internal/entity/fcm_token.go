package entity

import "time"

type FCMToken struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Token      string    `json:"token"`
	Platform   string    `json:"platform"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
