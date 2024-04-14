package dto

import "time"

type BannerHistoryOutput struct {
	CreatedAt time.Time      `json:"created_at"`
	Content   map[string]any `json:"content"`
	ID        int64          `json:"history_id"`
	BannerID  int64          `json:"banner_id"`
}
