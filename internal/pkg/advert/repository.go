package advert

type AdvertRepository interface {
	GetAdverts(page int, sortPrice, sortDate string, limit int) ([]Advert, error)
	GetAdvert(id int) (Advert, error)
	CreateAdvert(adv Advert) (int, error)
}
