package assetmodel

type Filter struct {
	OwnerId  int    `json:"-" form:"-"`
	WalletId int    `json:"-" form:"-"`
	Type     string `json:"-" form:"type"`
}
