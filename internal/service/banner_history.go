package service

import (
	"database/sql"

	"avito/internal/dto"
	"avito/internal/model"
	"avito/internal/repository"
)

type BannerHistory interface {
	Get(input *dto.GetBannerHistoryInput) ([]model.BannerHistory, error)
	Apply(input *dto.ApplyBannerHistoryInput) error
}

type bannerHistory struct {
	db                *sql.DB
	bannerRepository  repository.Banner
	historyRepository repository.BannerHistory
}

func NewBannerHistory(
	db *sql.DB,
	bannerRepository repository.Banner,
	historyRepository repository.BannerHistory,
) BannerHistory {
	return &bannerHistory{
		db:                db,
		bannerRepository:  bannerRepository,
		historyRepository: historyRepository,
	}
}

func (s *bannerHistory) Get(input *dto.GetBannerHistoryInput) ([]model.BannerHistory, error) {
	return s.historyRepository.FindByBannerID(input.BannerID)
}

func (s *bannerHistory) Apply(input *dto.ApplyBannerHistoryInput) error {
	err := repository.Transaction(s.db, func(tx *sql.Tx) error {
		bh, err := s.historyRepository.GetByID(tx, input.HistoryID)
		if err != nil {
			return err
		}

		b, err := s.bannerRepository.GetByID(tx, bh.BannerID)
		if err != nil {
			return err
		}

		if bh.Content == b.Content {
			return nil
		}

		if err = s.historyRepository.Create(tx, bh.BannerID, b.Content); err != nil {
			return err
		}

		if err = s.bannerRepository.UpdateContent(tx, bh.BannerID, bh.Content); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
