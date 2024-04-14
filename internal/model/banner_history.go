package model

import "time"

type BannerHistory struct {
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
	ID        int64     `json:"banner_history_id"`
	BannerID  int64     `json:"banner_id"`
}
