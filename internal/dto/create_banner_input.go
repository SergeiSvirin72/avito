package dto

type CreateBannerInput struct {
	Content   map[string]any `json:"content" binding:"required"`
	TagIDs    []int64        `json:"tag_ids" binding:"required,min=1,dive,gt=0"`
	FeatureID int64          `json:"feature_id" binding:"required,gt=0"`
	IsActive  bool           `json:"is_active"`
}
