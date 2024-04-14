package model

import "time"

type Banner struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`
	Tags      []Tag     `json:"tags"`
	ID        int64     `json:"id"`
	FeatureID int64     `json:"feature_id"`
	IsActive  bool      `json:"is_active"`
}

func GetBannerIDs(banners []Banner) []int64 {
	bannerIDs := make([]int64, len(banners))

	for i := range bannerIDs {
		bannerIDs[i] = banners[i].ID
	}

	return bannerIDs
}
