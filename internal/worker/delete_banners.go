package worker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"avito/internal/model"
	"avito/internal/repository"
)

const (
	TypeDeleteBanners   = "delete_banners"
	PeriodDeleteBanners = time.Second * 5
)

type DeleteBanners struct {
	db               *sql.DB
	logger           *slog.Logger
	jobRepository    repository.DeleteBannersJob
	bannerRepository repository.Banner
}

func NewDeleteBanners(
	db *sql.DB,
	logger *slog.Logger,
	jobRepository repository.DeleteBannersJob,
	bannerRepository repository.Banner,
) *DeleteBanners {
	return &DeleteBanners{
		db:               db,
		logger:           logger.With(slog.String("worker", TypeDeleteBanners)),
		jobRepository:    jobRepository,
		bannerRepository: bannerRepository,
	}
}

func (w *DeleteBanners) Run(ctx context.Context) {
	w.logger.Debug("Started")

	var (
		period time.Duration
		err    error
	)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(period):
			break
		}

		period, err = w.handle()
		if err != nil {
			w.logger.Error(err.Error())
			return
		}
	}
}

func (w *DeleteBanners) handle() (time.Duration, error) {
	task, err := w.jobRepository.GetFirst()
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return PeriodDeleteBanners, nil
	}

	if err != nil {
		return 0, err
	}

	switch {
	case task.FeatureID.Valid:
		if err = w.deleteByFeatureID(task); err != nil {
			return 0, err
		}
	case task.TagID.Valid:
		if err = w.deleteByTagID(task); err != nil {
			return 0, err
		}
	}

	w.logger.Debug(fmt.Sprintf("Done: featureID - %d, tagID - %d", task.FeatureID.Int64, task.TagID.Int64))

	return 0, nil
}

func (w *DeleteBanners) deleteByFeatureID(job *model.DeleteBannersJob) error {
	err := repository.Transaction(w.db, func(tx *sql.Tx) error {
		if err := w.bannerRepository.DeleteByFeatureID(tx, job.FeatureID.Int64); err != nil {
			return err
		}

		if err := w.jobRepository.DeleteByID(tx, job.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (w *DeleteBanners) deleteByTagID(job *model.DeleteBannersJob) error {
	err := repository.Transaction(w.db, func(tx *sql.Tx) error {
		if err := w.bannerRepository.DeleteByTagID(tx, job.TagID.Int64); err != nil {
			return err
		}

		if err := w.jobRepository.DeleteByID(tx, job.ID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
