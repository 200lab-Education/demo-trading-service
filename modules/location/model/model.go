package model

import "trading-service/common"

type Country struct {
	common.SQLModel
	Status int    `json:"status" gorm:"column:status;"`
	Name   string `json:"name" gorm:"column:name;"`
	Code   string `json:"code" gorm:"column:code;"`
	Cities []City `json:"cities"`
}

func (Country) TableName() string { return "countries" }

func (c *Country) Mask() {
	c.SQLModel.Mask(common.DbTypeCountry)

	if v := c.Cities; v != nil {
		for i := range v {
			v[i].Mask()
		}
	}
}

type City struct {
	common.SQLModel
	Status    int      `json:"status" gorm:"column:status;"`
	Name      string   `json:"name" gorm:"column:name;"`
	CountryId int      `json:"-" gorm:"column:country_id;"`
	Country   *Country `json:"country,omitempty" gorm:"foreignKey:CountryId;"`
}

func (City) TableName() string { return "cities" }

func (c *City) Mask() {
	c.SQLModel.Mask(common.DbTypeCity)

	if v := c.Country; v != nil {
		v.SQLModel.Mask(common.DbTypeCountry)
	}
}
