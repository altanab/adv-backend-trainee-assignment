package advert

type AdvertUsecase interface {
	GetAdverts(page int, sortPrice, sortDate string) ([]Advert, error)
	GetAdvert(id int, fields bool) (Advert, error)
	CreateAdvert(adv Advert) (int, error)
}
