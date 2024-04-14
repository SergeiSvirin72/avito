package dto

type ApplyBannerHistoryInput struct {
	HistoryID int64 `form:"history_id" binding:"required,gt=0"`
}
