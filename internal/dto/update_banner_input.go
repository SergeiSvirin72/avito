package dto

type UpdateBannerInput struct {
	Content   map[string]any `json:"content" validate:"required"`
	TagIDs    []int64        `json:"tag_ids" validate:"required,min=1,dive,gt=0"`
	ID        int64          `uri:"id" validate:"required,gt=0"`
	FeatureID int64          `json:"feature_id" validate:"required,gt=0"`
	IsActive  bool           `json:"is_active"`
}
