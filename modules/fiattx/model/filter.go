package model

type Filter struct {
	UserId      int    `json:"-" form:"-"`
	AssetId     int    `json:"-" form:"-"`
	Type        string `json:"type" form:"type"`
	Status      string `json:"status" form:"status"`
	PaymentType string `json:"payment_type" form:"payment_type"`
}
