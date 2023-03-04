package common

import (
	"time"
)

type SQLModel struct {
	Id     int  `json:"-" gorm:"column:id;"`
	FakeId *UID `json:"id" gorm:"-"`
	//Status    string     `json:"status" gorm:"column:status;"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;"`
}

func (m *SQLModel) Mask(dbType int) {
	if m == nil {
		return
	}

	uid := NewUID(uint32(m.Id), dbType, 1)
	m.FakeId = &uid
}

func NewSQLModel() SQLModel {
	now := time.Now().UTC()

	return SQLModel{
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}
