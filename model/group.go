package model

import (
	"time"
)

// Group : Model with injected fields `ID`, `CreatedAt`, `UpdatedAt`
type Group struct {
	GroupID   uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	GroupName string     `gorm:"type:varchar(100);not null"`
	GroupURL  string     `gorm:"type:varchar(180)"`
	GroupLogo string
	Verified  bool
	Active    bool
	Handle    string `gorm:"type:varchar(50);not null"`
	Ibans     []Iban `gorm:"polymorphic:Owner;"`
}
