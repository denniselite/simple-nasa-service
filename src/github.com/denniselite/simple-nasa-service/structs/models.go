package structs

import (
	"github.com/jinzhu/gorm"
	"time"
)

type NEO struct {
	gorm.Model
	Reference string
	Name string
	IsHazardous bool
}

type NEOData struct {
	gorm.Model
	Reference string
	Date time.Time
	Speed float64
}