package dto

type GetBannerInput struct {
	FeatureID       int64 `form:"feature_id" binding:"required,gt=0"`
	TagID           int64 `form:"tag_id" binding:"required,gt=0"`
	UseLastRevision bool  `form:"use_last_revision"`
}
