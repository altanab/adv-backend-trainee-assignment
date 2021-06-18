package repository

import (
	"adv/internal/pkg/advert"
	"errors"
	"gorm.io/gorm"
)

type AdvertRep struct {
	DB *gorm.DB
}

func (rep *AdvertRep) GetAdverts(page int, sortPrice, sortDate string, limit int) ([]advert.Advert, error) {
	adverts := make([]advert.Advert, 0)
	//если сортировка отсутствует, по умолчанию сортируется от "новых" к "старым"
	var orderQuery string
	switch {
	case sortPrice != "" && sortDate != "":
		orderQuery = "price " + sortPrice + ", created " + sortDate
	case sortPrice != "":
		orderQuery = "price " + sortPrice
	case sortDate != "":
		orderQuery = "created " + sortDate
	default:
		orderQuery = "created DESC"
	}
	rep.DB.
		Table("adverts").
		Select("id, adv_name, photos[1] as photo, price").
		Order(orderQuery).
		Limit(limit).
		Offset((page - 1) * limit).
		Scan(&adverts)
	if err := rep.DB.Error; err != nil {
		return nil, errors.New("cannot get adverts database error: " + err.Error())
	}
	return adverts, nil
}

func (rep *AdvertRep) GetAdvert(id int) (advert.Advert, error) {
	var adv advert.Advert
	rep.DB.
		Table("adverts").
		Select("adv_name, photos, price, about").
		Where("id=?", id).
		Scan(&adv)
	if err := rep.DB.Error; err != nil {
		return advert.Advert{}, errors.New("cannot get advert database error: " + err.Error())
	}
	return adv, nil
}

func (rep *AdvertRep) CreateAdvert(adv advert.Advert) (int, error) {
	result := rep.DB.Select("adv_name", "about", "photos", "price").Create(&adv)
	if err := result.Error; err != nil {
		return 0, errors.New("cannot create advert database error: " + err.Error())
	}
	return adv.Id, nil
}