package usecase

import (
	"adv/internal/pkg/advert"
	"errors"
	"strings"
)

type AdvertUC struct {
	AdvertRep advert.AdvertRepository
	PaginationLimit int `json:"pagination_limit"`
	MaxNumPhotos int `json:"max_num_photos"`
	MaxCharName int `json:"max_char_name"`
	MaxCharAbout int `json:"max_char_about"`
}

func (uc *AdvertUC) GetAdverts(page int, sortPrice, sortDate string) ([]advert.Advert, error) {
	sortPrice = strings.ToUpper(sortPrice)
	sortDate = strings.ToUpper(sortDate)
	if  sortPrice != "ASC" && sortPrice != "DESC" {
		sortPrice = ""
	}
	if  sortDate != "ASC" && sortDate != "DESC" {
		sortDate = ""
	}
	return uc.AdvertRep.GetAdverts(page, sortPrice, sortDate, uc.PaginationLimit)
}

func (uc *AdvertUC) GetAdvert(id int, fields bool) (advert.Advert, error) {
	adv, err := uc.AdvertRep.GetAdvert(id)
	if err != nil {
		return advert.Advert{}, err
	}
	if !fields {
		adv.About = ""
		if adv.Photos != nil {
			adv.MainPhoto = adv.Photos[0]
			adv.Photos = nil
		}
	}
	return adv, nil
}

func (uc *AdvertUC) CreateAdvert(adv advert.Advert) (int, error) {
	if len(adv.Photos) > uc.MaxNumPhotos ||
		len(adv.Name) > uc.MaxCharName ||
		len(adv.About) > uc.MaxCharAbout {
		return 0, errors.New("incorrect data")
	}
	return uc.AdvertRep.CreateAdvert(adv)
}