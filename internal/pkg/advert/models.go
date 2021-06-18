package advert

import (
	"github.com/lib/pq"
	"time"
)

type Advert struct {
	Id int `json:"id,omitempty"`
	Name string `json:"name" gorm:"column:adv_name"`
	About string `json:"about,omitempty"`
	Photos pq.StringArray `json:"photos,omitempty" gorm:"type:varchar[]"`
	Price int `json:"price"`
	Created time.Time `json:"-"`
	MainPhoto string `json:"photo,omitempty" gorm:"->;column:photo"`
}
