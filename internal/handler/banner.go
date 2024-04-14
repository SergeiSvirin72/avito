package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"avito/internal/dto"
	"avito/internal/presenter"
	"avito/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Banner struct {
	validate        *validator.Validate
	bannerService   service.Banner
	bannerPresenter *presenter.Banner
}

func NewBanner(validate *validator.Validate, bannerService service.Banner, bannerPresenter *presenter.Banner) *Banner {
	return &Banner{
		validate:        validate,
		bannerService:   bannerService,
		bannerPresenter: bannerPresenter,
	}
}

func (h *Banner) Get(c *gin.Context) {
	var input dto.GetBannersInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	banners, err := h.bannerService.GetBanners(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output, err := h.bannerPresenter.Present(banners)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (h *Banner) Create(c *gin.Context) {
	var input dto.CreateBannerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.bannerService.Create(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, h.bannerPresenter.PresentID(id))
}

func (h *Banner) Update(c *gin.Context) {
	var input dto.UpdateBannerInput
	if err := c.ShouldBindUri(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.bannerService.Update(&input)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Banner) Delete(c *gin.Context) {
	var input dto.DeleteBannerInput
	if err := c.ShouldBindUri(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := h.bannerService.Delete(&input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Banner) DeleteMany(c *gin.Context) {
	var input dto.DeleteBannersInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.bannerService.CreateDeleteJob(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{"message": "banners will be deleted"})
}
