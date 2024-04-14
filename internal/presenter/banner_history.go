package presenter

import (
	"encoding/json"
	"fmt"

	"avito/internal/dto"
	"avito/internal/model"
)

type BannerHistory struct {
}

func NewBannerHistory() *BannerHistory {
	return &BannerHistory{}
}

func (p *BannerHistory) Present(bannersHistory []model.BannerHistory) ([]dto.BannerHistoryOutput, error) {
	output := make([]dto.BannerHistoryOutput, len(bannersHistory))

	for i := range bannersHistory {
		output[i] = dto.BannerHistoryOutput{
			ID:        bannersHistory[i].ID,
			BannerID:  bannersHistory[i].BannerID,
			CreatedAt: bannersHistory[i].CreatedAt,
		}

		if err := json.Unmarshal([]byte(bannersHistory[i].Content), &output[i].Content); err != nil {
			return nil, fmt.Errorf("bannerHistoryMapper.Present, unmarshall error: %w", err)
		}
	}

	return output, nil
}
