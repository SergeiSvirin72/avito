package dto

type DeleteBannerInput struct {
	ID int64 `uri:"id" binding:"required,gt=0"`
}
