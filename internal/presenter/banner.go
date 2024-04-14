package presenter

import (
	"encoding/json"
	"fmt"

	"avito/internal/dto"
	"avito/internal/model"
)

type Banner struct {
}

func NewBanner() *Banner {
	return &Banner{}
}

func (p *Banner) Present(banners []model.Banner) ([]dto.BannerOutput, error) {
	output := make([]dto.BannerOutput, len(banners))

	for i := range banners {
		tagIDs := make([]int64, len(banners[i].Tags))

		for j := range banners[i].Tags {
			tagIDs[j] = banners[i].Tags[j].ID
		}

		output[i] = dto.BannerOutput{
			ID:        banners[i].ID,
			FeatureID: banners[i].FeatureID,
			IsActive:  banners[i].IsActive,
			CreatedAt: banners[i].CreatedAt,
			UpdatedAt: banners[i].UpdatedAt,
			TagIDs:    tagIDs,
		}

		if err := json.Unmarshal([]byte(banners[i].Content), &output[i].Content); err != nil {
			return nil, fmt.Errorf("bannerMapper.Present, unmarshall error: %w", err)
		}
	}

	return output, nil
}

func (p *Banner) PresentID(id int64) *dto.BannerIDOutput {
	return &dto.BannerIDOutput{
		BannerID: id,
	}
}
