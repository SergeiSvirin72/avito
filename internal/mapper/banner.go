package mapper

import (
	"avito/internal/model"
)

func MapBannerTags(banners []model.Banner, bannerTags []model.BannerTag) []model.Banner {
	tagsMap := make(map[int64][]model.Tag)

	for i := range bannerTags {
		row := bannerTags[i]

		_, ok := tagsMap[row.BannerID]
		if !ok {
			tagsMap[row.BannerID] = make([]model.Tag, 0)
		}

		tagsMap[row.BannerID] = append(tagsMap[row.BannerID], model.Tag{ID: row.TagID})
	}

	for i := range banners {
		banners[i].Tags = tagsMap[banners[i].ID]
	}

	return banners
}
