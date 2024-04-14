package repository

import (
	"database/sql"
	"fmt"

	"avito/internal/model"
)

type DeleteBannersJob interface {
	Create(*model.DeleteBannersJob) error
	DeleteByID(tx *sql.Tx, id int64) error
	GetFirst() (*model.DeleteBannersJob, error)
}

type deleteBannersJob struct {
	db *sql.DB
}

func NewDeleteBannersJob(db *sql.DB) DeleteBannersJob {
	return &deleteBannersJob{
		db: db,
	}
}

func (r *deleteBannersJob) Create(task *model.DeleteBannersJob) error {
	_, err := r.db.Exec(
		"INSERT INTO delete_banners_job (feature_id, tag_id) VALUES ($1, $2)",
		task.FeatureID,
		task.TagID,
	)
	if err != nil {
		return fmt.Errorf("deleteBannersJobRepository.Create: %w", err)
	}

	return nil
}

func (r *deleteBannersJob) DeleteByID(tx *sql.Tx, id int64) error {
	_, err := tx.Exec("DELETE FROM delete_banners_job WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("deleteBannersJobRepository.DeleteByID: %w", err)
	}

	return nil
}

func (r *deleteBannersJob) GetFirst() (*model.DeleteBannersJob, error) {
	var t model.DeleteBannersJob

	row := r.db.QueryRow("SELECT * FROM delete_banners_job ORDER BY created_at LIMIT 1")
	if err := row.Scan(&t.ID, &t.FeatureID, &t.TagID, &t.CreatedAt); err != nil {
		return nil, fmt.Errorf("deleteBannersJobRepository.GetFirst: %w", err)
	}

	return &t, nil
}
