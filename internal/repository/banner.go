package repository

import (
	"database/sql"
	"fmt"

	"avito/internal/mapper"
	"avito/internal/model"

	"github.com/lib/pq"
)

type Banner interface {
	Create(tx *sql.Tx, b *model.Banner) (int64, error)
	Update(tx *sql.Tx, b *model.Banner) error
	UpdateContent(tx *sql.Tx, id int64, content string) error
	DeleteByID(id int64) (int64, error)
	DeleteByFeatureID(tx *sql.Tx, featureID int64) error
	DeleteByTagID(tx *sql.Tx, tagID int64) error
	GetByID(tx *sql.Tx, id int64) (*model.Banner, error)
	GetByFeatureIDAndTagID(featureID, tagID int64) (*model.Banner, error)
	FindByFeatureIDAndTagID(featureID, tagID, limit, offset int64) ([]model.Banner, error)
	FindByFeatureID(featureID, limit, offset int64) ([]model.Banner, error)
	FindByTagID(tagID, limit, offset int64) ([]model.Banner, error)
	Find(limit, offset int64) ([]model.Banner, error)
}

type banner struct {
	db *sql.DB
}

func NewBanner(db *sql.DB) Banner {
	return &banner{
		db: db,
	}
}

func (r *banner) Create(tx *sql.Tx, b *model.Banner) (int64, error) {
	var id int64

	row := tx.QueryRow(
		"INSERT INTO banners (feature_id, content, is_active) VALUES ($1, $2, $3) RETURNING id",
		b.FeatureID,
		b.Content,
		b.IsActive,
	)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("bannerRepository.Create: %w", err)
	}

	return id, nil
}

func (r *banner) Update(tx *sql.Tx, b *model.Banner) error {
	_, err := tx.Exec(
		"UPDATE banners SET feature_id = $1, content = $2, is_active = $3, updated_at = now() WHERE id = $4",
		b.FeatureID,
		b.Content,
		b.IsActive,
		b.ID,
	)
	if err != nil {
		return fmt.Errorf("bannerRepository.Update: %w", err)
	}

	return nil
}

func (r *banner) UpdateContent(tx *sql.Tx, id int64, content string) error {
	_, err := tx.Exec("UPDATE banners SET content = $1, updated_at = now() WHERE id = $2", content, id)
	if err != nil {
		return fmt.Errorf("bannerRepository.UpdateContent: %w", err)
	}

	return nil
}

func (r *banner) DeleteByID(id int64) (int64, error) {
	res, err := r.db.Exec("DELETE FROM banners WHERE id = $1", id)
	if err != nil {
		return 0, fmt.Errorf("bannerRepository.DeleteByID %d: %w", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("bannerRepository.DeleteByID %d: %w", id, err)
	}

	return rowsAffected, nil
}

func (r *banner) DeleteByFeatureID(tx *sql.Tx, featureID int64) error {
	_, err := tx.Exec("DELETE FROM banners WHERE feature_id = $1", featureID)
	if err != nil {
		return fmt.Errorf("bannerRepository.DeleteByFeatureID %d: %w", featureID, err)
	}

	return nil
}

func (r *banner) DeleteByTagID(tx *sql.Tx, tagID int64) error {
	_, err := tx.Exec(`DELETE FROM banners b USING banner_tag bt WHERE b.id = bt.banner_id AND bt.tag_id = $1;`, tagID)
	if err != nil {
		return fmt.Errorf("bannerRepository.DeleteByTagID %d: %w", tagID, err)
	}

	return nil
}

func (r *banner) GetByID(tx *sql.Tx, id int64) (*model.Banner, error) {
	var b model.Banner

	row := tx.QueryRow(`SELECT b.content FROM banners b WHERE b.id = $1`, id)
	if err := row.Scan(&b.Content); err != nil {
		return nil, fmt.Errorf("bannerRepository.GetByID %d: %w", id, err)
	}

	return &b, nil
}

func (r *banner) GetByFeatureIDAndTagID(featureID, tagID int64) (*model.Banner, error) {
	var b model.Banner

	row := r.db.QueryRow(`
		SELECT b.content, b.is_active FROM banners b 
		JOIN banner_tag bt ON b.id = bt.banner_id
		WHERE b.feature_id = $1 AND bt.tag_id = $2
	`, featureID, tagID)

	if err := row.Scan(&b.Content, &b.IsActive); err != nil {
		return nil, fmt.Errorf("bannerRepository.GetByFeatureIDAndTagID %d %d: %w", featureID, tagID, err)
	}

	return &b, nil
}

func (r *banner) FindByFeatureIDAndTagID(featureID, tagID, limit, offset int64) ([]model.Banner, error) {
	banners, err := r.find(`
		SELECT DISTINCT b.* FROM banners b
		JOIN banner_tag bt ON b.id = bt.banner_id
		WHERE b.feature_id = $1 AND bt.tag_id = $2
		ORDER BY b.id
		LIMIT $3 OFFSET $4
	`, featureID, tagID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf(
			"bannerRepository.FindByFeatureIDAndTagID %d %d %d %d: %w",
			featureID, tagID, limit, offset, err,
		)
	}

	return banners, nil
}

func (r *banner) FindByFeatureID(featureID, limit, offset int64) ([]model.Banner, error) {
	banners, err := r.find(`
		SELECT DISTINCT b.* FROM banners b
		JOIN banner_tag bt ON b.id = bt.banner_id
		WHERE b.feature_id = $1
		ORDER BY b.id
		LIMIT $2 OFFSET $3
	`, featureID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("bannerRepository.FindByFeatureID %d %d %d: %w", featureID, limit, offset, err)
	}

	return banners, nil
}

func (r *banner) FindByTagID(tagID, limit, offset int64) ([]model.Banner, error) {
	banners, err := r.find(`
		SELECT DISTINCT b.* FROM banners b
		JOIN banner_tag bt ON b.id = bt.banner_id
		WHERE bt.tag_id = $1
		ORDER BY b.id
		LIMIT $2 OFFSET $3
	`, tagID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("bannerRepository.FindByTagID %d %d %d: %w", tagID, limit, offset, err)
	}

	return banners, nil
}

func (r *banner) Find(limit, offset int64) ([]model.Banner, error) {
	banners, err := r.find(`
		SELECT DISTINCT b.* FROM banners b
		JOIN banner_tag bt ON b.id = bt.banner_id
		ORDER BY b.id
		LIMIT $1 OFFSET $2
	`, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("bannerRepository.Find %d %d: %w", limit, offset, err)
	}

	return banners, nil
}

func (r *banner) find(query string, args ...any) ([]model.Banner, error) {
	bRows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer bRows.Close()

	banners := make([]model.Banner, 0)

	for bRows.Next() {
		var b model.Banner
		if err = bRows.Scan(&b.ID, &b.FeatureID, &b.Content, &b.IsActive, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}

		banners = append(banners, b)
	}

	if err = bRows.Err(); err != nil {
		return nil, err
	}

	if len(banners) == 0 {
		return banners, nil
	}

	btRows, err := r.db.Query(
		"SELECT bt.* FROM banner_tag bt WHERE bt.banner_id = any($1)",
		pq.Array(model.GetBannerIDs(banners)),
	)
	if err != nil {
		return nil, err
	}

	defer btRows.Close()

	bannerTags := make([]model.BannerTag, 0)

	for btRows.Next() {
		var bt model.BannerTag
		if err = btRows.Scan(&bt.BannerID, &bt.TagID); err != nil {
			return nil, err
		}

		bannerTags = append(bannerTags, bt)
	}

	if err = btRows.Err(); err != nil {
		return nil, err
	}

	return mapper.MapBannerTags(banners, bannerTags), nil
}
