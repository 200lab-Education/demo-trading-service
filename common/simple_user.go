package common

type SimpleUser struct {
	SQLModel  `json:",inline"`
	LastName  string `json:"last_name" gorm:"column:last_name;"`
	FirstName string `json:"first_name" gorm:"column:first_name;"`
	Role      string `json:"role" gorm:"column:role;"`
	Active    int    `json:"active" gorm:"column:active;"`
}

func (SimpleUser) TableName() string {
	return "app_users"
}
