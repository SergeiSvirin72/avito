package model

import (
	"database/sql"
	"time"
)

type DeleteBannersJob struct {
	CreatedAt time.Time
	ID        int64
	FeatureID sql.NullInt64
	TagID     sql.NullInt64
}

func (m *DeleteBannersJob) SetFeatureID(featureID *int64) {
	if featureID == nil {
		m.FeatureID = sql.NullInt64{
			Int64: 0,
			Valid: false,
		}

		return
	}

	m.FeatureID = sql.NullInt64{
		Int64: *featureID,
		Valid: true,
	}
}

func (m *DeleteBannersJob) SetTagID(tagID *int64) {
	if tagID == nil {
		m.TagID = sql.NullInt64{
			Int64: 0,
			Valid: false,
		}

		return
	}

	m.TagID = sql.NullInt64{
		Int64: *tagID,
		Valid: true,
	}
}
