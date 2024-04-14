package dto

import "time"

type BannerOutput struct {
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Content   map[string]any `json:"content"`
	TagIDs    []int64        `json:"tag_ids"`
	ID        int64          `json:"banner_id"`
	FeatureID int64          `json:"feature_id"`
	IsActive  bool           `json:"is_active"`
}
