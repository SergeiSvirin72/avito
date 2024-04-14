package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"avito/internal/dto"
	"avito/internal/middleware"
	"avito/internal/model"
	"avito/internal/presenter"
	"avito/internal/service"

	"github.com/gin-gonic/gin"
)

type UserBanner struct {
	bannerService   service.Banner
	bannerPresenter *presenter.Banner
}

func NewUserBanner(bannerService service.Banner, bannerPresenter *presenter.Banner) *UserBanner {
	return &UserBanner{
		bannerService:   bannerService,
		bannerPresenter: bannerPresenter,
	}
}

func (h *UserBanner) Get(c *gin.Context) {
	var input dto.GetBannerInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := c.Get(middleware.KeyUser)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
		return
	}

	b, err := h.bannerService.GetBanner(c.Request.Context(), &input, user.(*model.User))
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err != nil && errors.Is(err, service.ErrForbidden) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, gin.MIMEJSON, []byte(b.Content))
}
