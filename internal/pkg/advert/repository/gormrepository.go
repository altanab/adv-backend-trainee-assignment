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
	var orderQuery string
	if sortPrice != "" {
		orderQuery += "price " + sortPrice + " "
	}
	if sortDate != "" {
		orderQuery += "created " + sortDate
	}

	//использование limit и offset не самый консинстентный и эффективный, но самый простой вариант

	rep.DB.
		Table("adverts").
		Select("id, name, photos[1], price").
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
		Select("name, photos, price, about").
		Where("id=?", id).
		Scan(&adv)
	if err := rep.DB.Error; err != nil {
		return advert.Advert{}, errors.New("cannot get advert database error: " + err.Error())
	}
	return adv, nil
}

func (rep *AdvertRep) CreateAdvert(adv advert.Advert) (int, error) {
	result := rep.DB.Omit("name", "about", "photos", "price").Create(&adv)
	if err := result.Error; err != nil {
		return 0, errors.New("cannot create advert database error: " + err.Error())
	}
	return adv.Id, nil
}