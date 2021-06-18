package advert

import "time"

type Advert struct {
	Id int `json:"id"`
	Name string `json:"name"`
	About string `json:"about,omitempty"`
	Photos []string `json:"photos,omitempty"`
	Price int `json:"price"`
	Created time.Time `json:"-"`
	MainPhoto string `json:"photo,omitempty"`
}