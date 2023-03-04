package usermodel

import (
	"errors"
	"time"
	"trading-service/common"
)

const (
	EntityName = "User"
	TbName     = "app_users"
)

type User struct {
	common.SQLModel     `json:",inline"`
	Email               string     `json:"email" gorm:"column:email;"`
	Username            string     `json:"username" gorm:"column:username;"`
	Password            string     `json:"-" gorm:"column:password;"`
	FirstName           string     `json:"first_name" gorm:"column:first_name;"`
	LastName            string     `json:"last_name" gorm:"column:last_name;"`
	Salt                string     `json:"-" gorm:"column:salt;"`
	WalletAddress       string     `json:"eth_address" gorm:"column:eth_address;"`
	EthWalletVerified   bool       `json:"eth_address_verified" gorm:"column:eth_address_verified;"`
	EthWalletVerifiedAt *time.Time `json:"eth_address_verified_at" gorm:"column:eth_address_verified_at;"`
	IsSetPassword       bool       `json:"is_set_password" gorm:"column:is_set_password;"`
	Active              int        `json:"active" gorm:"column:active;"`
	Role                string     `json:"role" gorm:"column:role;"`
}

type UserUpdate struct {
	FirstName *string `json:"first_name" gorm:"column:first_name;"`
	LastName  *string `json:"last_name" gorm:"column:last_name;"`
	Active    *int    `json:"-" gorm:"column:active;"`
}

type UserCreate struct {
	common.SQLModel `json:",inline"`
	Username        string `json:"username" gorm:"column:username;"`
	Password        string `json:"password" gorm:"column:password;"`
	LastName        string `json:"last_name" gorm:"column:last_name;"`
	FirstName       string `json:"first_name" gorm:"column:first_name;"`
	Role            string `json:"-" gorm:"column:role;"`
	Salt            string `json:"-" gorm:"column:salt;"`
}

func (UserUpdate) TableName() string { return User{}.TableName() }
func (UserCreate) TableName() string { return User{}.TableName() }

func (u *User) GetUserId() int {
	return u.Id
}

func (u *User) GetRole() string {
	return u.Role
}

func (u *User) GetWalletAddress() string {
	return u.WalletAddress
}

func (User) TableName() string {
	return TbName
}

type UserLogin struct {
	Username string `json:"username" form:"username" gorm:"column:username;"`
	Password string `json:"password" form:"password" gorm:"column:password;"`
}

func (UserLogin) TableName() string {
	return User{}.TableName()
}

var (
	ErrUsernameOrPasswordInvalid = common.NewCustomError(
		errors.New("username or password invalid"),
		"username or password invalid",
		"ErrUsernameOrPasswordInvalid",
	)

	ErrUsernameExisted = common.NewCustomError(
		errors.New("username has already existed"),
		"username has already existed",
		"ErrUsernameExisted",
	)
)
