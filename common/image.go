package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Image struct {
	Filename  string `json:"filename" gorm:"column:filename;"`
	Width     int    `json:"width" gorm:"column:width;"`
	Height    int    `json:"height" gorm:"column:height;"`
	CloudName string `json:"cloud" gorm:"column:cloud;"`
	Url       string `json:"url" gorm:"-"`
}

func (Image) TableName() string { return "images" }

func (j *Image) Fulfill(domain string) {
	if j == nil {
		return
	}

	j.Url = fmt.Sprintf("%s/%s", domain, j.Filename)
}

func (j *Image) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var img Image
	if err := json.Unmarshal(bytes, &img); err != nil {
		return err
	}

	*j = img
	return nil
}

// Value return json value, implement driver.Valuer interface
func (j *Image) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

type Images []Image

func (j *Images) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var img []Image
	if err := json.Unmarshal(bytes, &img); err != nil {
		return err
	}

	*j = img
	return nil
}

// Value return json value, implement driver.Valuer interface
func (j *Images) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
