package dto

type DeleteBannersInput struct {
	FeatureID *int64 `form:"feature_id" binding:"required_without=TagID,excluded_with=TagID,omitempty,gt=0"`
	TagID     *int64 `form:"tag_id" binding:"required_without=FeatureID,excluded_with=FeatureID,omitempty,gt=0"`
}
