package dto

type GetBannerHistoryInput struct {
	BannerID int64 `form:"banner_id" binding:"required,gt=0"`
}
