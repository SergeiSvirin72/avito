package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"avito/internal/dto"
	"avito/internal/model"
	"avito/internal/repository"
)

const (
	BannerCacheKey = "banner_feature_id_%d_tag_id_%d"
	BannerTTL      = time.Minute * 5
)

var ErrForbidden = errors.New("forbidden")

type Banner interface {
	GetBanner(ctx context.Context, input *dto.GetBannerInput, user *model.User) (*model.Banner, error)
	GetBanners(input *dto.GetBannersInput) ([]model.Banner, error)
	Create(input *dto.CreateBannerInput) (int64, error)
	Update(input *dto.UpdateBannerInput) error
	Delete(input *dto.DeleteBannerInput) (int64, error)
	CreateDeleteJob(input *dto.DeleteBannersInput) error
}

type banner struct {
	db                  *sql.DB
	cacheService        Cache
	bannerRepository    repository.Banner
	bannerTagRepository repository.BannerTag
	historyRepository   repository.BannerHistory
	jobRepository       repository.DeleteBannersJob
}

func NewBanner(
	db *sql.DB,
	cacheService Cache,
	bannerRepository repository.Banner,
	bannerTagRepository repository.BannerTag,
	historyRepository repository.BannerHistory,
	jobRepository repository.DeleteBannersJob,
) Banner {
	return &banner{
		db:                  db,
		cacheService:        cacheService,
		bannerRepository:    bannerRepository,
		bannerTagRepository: bannerTagRepository,
		historyRepository:   historyRepository,
		jobRepository:       jobRepository,
	}
}

func (s *banner) GetBanner(ctx context.Context, input *dto.GetBannerInput, user *model.User) (*model.Banner, error) {
	if !input.UseLastRevision {
		b, err := s.getFromCache(ctx, input.FeatureID, input.TagID)
		if err != nil && !errors.Is(err, ErrKeyNotExist) {
			return nil, err
		}

		if b != nil && !b.IsActive && !user.IsAdmin() {
			return nil, ErrForbidden
		}

		if b != nil {
			return b, nil
		}
	}

	b, err := s.bannerRepository.GetByFeatureIDAndTagID(input.FeatureID, input.TagID)
	if err != nil {
		return nil, err
	}

	val, err := json.Marshal(&b)
	if err != nil {
		return nil, fmt.Errorf("bannerService.GetBanner, marshall error: %w", err)
	}

	if err = s.cacheService.Set(ctx, s.getCacheKey(input.FeatureID, input.TagID), val, BannerTTL); err != nil {
		return nil, err
	}

	if !b.IsActive && !user.IsAdmin() {
		return nil, ErrForbidden
	}

	return b, nil
}

func (s *banner) GetBanners(input *dto.GetBannersInput) ([]model.Banner, error) {
	var (
		banners []model.Banner
		err     error
	)

	// т. к. нет query builder и метод в репозитории получился бы сложным - разделил на 4 разных
	switch {
	case input.FeatureID == nil && input.TagID == nil:
		banners, err = s.bannerRepository.Find(*input.Limit, *input.Offset)
	case input.FeatureID == nil && input.TagID != nil:
		banners, err = s.bannerRepository.FindByTagID(*input.TagID, *input.Limit, *input.Offset)
	case input.FeatureID != nil && input.TagID == nil:
		banners, err = s.bannerRepository.FindByFeatureID(*input.FeatureID, *input.Limit, *input.Offset)
	case input.FeatureID != nil && input.TagID != nil:
		banners, err = s.bannerRepository.FindByFeatureIDAndTagID(
			*input.FeatureID,
			*input.TagID,
			*input.Limit,
			*input.Offset,
		)
	}

	if err != nil {
		return nil, err
	}

	return banners, nil
}

func (s *banner) Create(input *dto.CreateBannerInput) (int64, error) {
	content, err := json.Marshal(input.Content)
	if err != nil {
		return 0, fmt.Errorf("bannerService.Create, marshall error: %w", err)
	}

	b := model.Banner{
		FeatureID: input.FeatureID,
		Content:   string(content),
		IsActive:  input.IsActive,
	}

	var id int64

	err = repository.Transaction(s.db, func(tx *sql.Tx) error {
		id, err = s.bannerRepository.Create(tx, &b)
		if err != nil {
			return err
		}

		if err = s.bannerTagRepository.Create(tx, id, input.TagIDs); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *banner) Update(input *dto.UpdateBannerInput) error {
	content, err := json.Marshal(input.Content)
	if err != nil {
		return fmt.Errorf("bannerService.Update, marshall error: %w", err)
	}

	b := model.Banner{
		ID:        input.ID,
		FeatureID: input.FeatureID,
		Content:   string(content),
		IsActive:  input.IsActive,
	}

	var oldB *model.Banner

	err = repository.Transaction(s.db, func(tx *sql.Tx) error {
		oldB, err = s.bannerRepository.GetByID(tx, input.ID)
		if err != nil {
			return err
		}

		if err = s.bannerRepository.Update(tx, &b); err != nil {
			return err
		}

		if err = s.bannerTagRepository.DeleteByBannerID(tx, input.ID); err != nil {
			return err
		}

		if err = s.bannerTagRepository.Create(tx, b.ID, input.TagIDs); err != nil {
			return err
		}

		if b.Content == oldB.Content {
			return nil
		}

		if err = s.historyRepository.Create(tx, input.ID, oldB.Content); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *banner) Delete(input *dto.DeleteBannerInput) (int64, error) {
	count, err := s.bannerRepository.DeleteByID(input.ID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *banner) CreateDeleteJob(input *dto.DeleteBannersInput) error {
	var j model.DeleteBannersJob

	j.SetFeatureID(input.FeatureID)
	j.SetTagID(input.TagID)

	if err := s.jobRepository.Create(&j); err != nil {
		return err
	}

	return nil
}

func (s *banner) getCacheKey(featureID, tagID int64) string {
	return fmt.Sprintf(BannerCacheKey, featureID, tagID)
}

func (s *banner) getFromCache(ctx context.Context, featureID, tagID int64) (*model.Banner, error) {
	val, err := s.cacheService.Get(ctx, s.getCacheKey(featureID, tagID))
	if err != nil {
		return nil, err
	}

	var b model.Banner

	if err = json.Unmarshal([]byte(val), &b); err != nil {
		return nil, fmt.Errorf("bannerService.getFromCache, unmarshall error: %w", err)
	}

	return &b, nil
}
