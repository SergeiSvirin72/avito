package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

const placeholderOffset = 2

type BannerTag interface {
	Create(tx *sql.Tx, bannerID int64, tagIDs []int64) error
	DeleteByBannerID(tx *sql.Tx, bannerID int64) error
}

type bannerTag struct {
	db *sql.DB
}

func NewBannerTag(db *sql.DB) BannerTag {
	return &bannerTag{
		db: db,
	}
}

func (r *bannerTag) Create(tx *sql.Tx, bannerID int64, tagIDs []int64) error {
	args := make([]any, len(tagIDs)+1)
	args[0] = bannerID

	query := "INSERT INTO banner_tag (banner_id, tag_id) VALUES "

	for i := range tagIDs {
		args[i+1] = tagIDs[i]
		query += fmt.Sprintf("($1, $%d), ", i+placeholderOffset)
	}

	query = strings.TrimSuffix(query, ", ")

	_, err := tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("bannerTagRepository.Create: %w", err)
	}

	return nil
}

func (r *bannerTag) DeleteByBannerID(tx *sql.Tx, bannerID int64) error {
	_, err := tx.Exec("DELETE FROM banner_tag WHERE banner_id = $1", bannerID)
	if err != nil {
		return fmt.Errorf("bannerTagRepository.DeleteByBannerID: %w", err)
	}

	return nil
}
