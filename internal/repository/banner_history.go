package repository

import (
	"database/sql"
	"fmt"

	"avito/internal/model"
)

type BannerHistory interface {
	Create(tx *sql.Tx, bannerID int64, content string) error
	GetByID(tx *sql.Tx, id int64) (*model.BannerHistory, error)
	FindByBannerID(bannerID int64) ([]model.BannerHistory, error)
}

type bannerHistory struct {
	db *sql.DB
}

func NewBannerHistory(db *sql.DB) BannerHistory {
	return &bannerHistory{
		db: db,
	}
}

func (r *bannerHistory) Create(tx *sql.Tx, bannerID int64, content string) error {
	_, err := tx.Exec("INSERT INTO banner_history (banner_id, content) VALUES ($1, $2)", bannerID, content)
	if err != nil {
		return fmt.Errorf("bannerHistoryRepository.Create: %w", err)
	}

	return nil
}

func (r *bannerHistory) GetByID(tx *sql.Tx, id int64) (*model.BannerHistory, error) {
	var bh model.BannerHistory

	row := tx.QueryRow(`SELECT bh.banner_id, bh.content FROM banner_history bh WHERE bh.id = $1`, id)
	if err := row.Scan(&bh.BannerID, &bh.Content); err != nil {
		return nil, fmt.Errorf("bannerHistoryRepository.GetByID %d: %w", id, err)
	}

	return &bh, nil
}

func (r *bannerHistory) FindByBannerID(bannerID int64) ([]model.BannerHistory, error) {
	rows, err := r.db.Query(
		"SELECT * FROM banner_history WHERE banner_id = $1 ORDER BY created_at DESC LIMIT 3",
		bannerID,
	)
	if err != nil {
		return nil, fmt.Errorf("bannerHistoryRepository.FindByBannerID %d: %w", bannerID, err)
	}

	defer rows.Close()

	bannersHistory := make([]model.BannerHistory, 0)

	for rows.Next() {
		var bh model.BannerHistory
		if err = rows.Scan(&bh.ID, &bh.BannerID, &bh.Content, &bh.CreatedAt); err != nil {
			return nil, fmt.Errorf("bannerHistoryRepository.FindByBannerID %d: %w", bannerID, err)
		}

		bannersHistory = append(bannersHistory, bh)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("bannerHistoryRepository.FindByBannerID %d: %w", bannerID, err)
	}

	return bannersHistory, nil
}
