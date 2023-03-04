package model

type Filter struct {
	UserId          int  `json:"-" form:"-"`
	IsSystemAccount bool `json:"-" form:"-"`
}
