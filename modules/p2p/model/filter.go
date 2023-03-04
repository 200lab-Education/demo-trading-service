package model

type Filter struct {
	UserId       string `json:"-" form:"user_id"`
	OrderId      int    `json:"-" form:"-"`
	OfferAssetId string `json:"-" form:"offer_asset_id"`
	Type         string `json:"-" form:"type"`
	Status       string `json:"-" form:"status"`
	IsAdmin      bool   `json:"-" form:"-"`
	UserAccess   bool   `json:"-" form:"-"`
	RequesterId  int    `json:"-" form:"-"`
}
