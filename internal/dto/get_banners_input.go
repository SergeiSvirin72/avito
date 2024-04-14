package dto

type GetBannersInput struct {
	FeatureID *int64 `form:"feature_id" binding:"omitempty,gt=0"`
	TagID     *int64 `form:"tag_id" binding:"omitempty,gt=0"`
	Limit     *int64 `form:"limit,default=100" binding:"omitempty,gt=0,lte=100"`
	Offset    *int64 `form:"offset,default=0" binding:"omitempty,gte=0"`
}
