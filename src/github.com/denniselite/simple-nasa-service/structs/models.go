package structs

import (
	"time"
)

type NEO struct {
	Reference   string `gorm:"primary_key" sql:"not null;unique"`
	Name        string
	IsHazardous bool
	NEOData     []NEOData `gorm:"ForeignKey:Reference"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
}

type NEOData struct {
	Id        int `gorm:"primary_key"`
	Reference string `sql:"index"`
	Date      time.Time
	Speed     float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
