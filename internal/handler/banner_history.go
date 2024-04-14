package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"avito/internal/dto"
	"avito/internal/presenter"
	"avito/internal/service"

	"github.com/gin-gonic/gin"
)

type BannerHistory struct {
	historyService   service.BannerHistory
	historyPresenter *presenter.BannerHistory
}

func NewBannerHistory(historyService service.BannerHistory, historyPresenter *presenter.BannerHistory) *BannerHistory {
	return &BannerHistory{
		historyService:   historyService,
		historyPresenter: historyPresenter,
	}
}

func (h *BannerHistory) Get(c *gin.Context) {
	var input dto.GetBannerHistoryInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	history, err := h.historyService.Get(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output, err := h.historyPresenter.Present(history)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *BannerHistory) Apply(c *gin.Context) {
	var input dto.ApplyBannerHistoryInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.historyService.Apply(&input)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
