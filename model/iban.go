package model

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Iban : Model with injected fields `ID`, `CreatedAt`, `UpdatedAt`
type Iban struct {
	IbanID      uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	Text        string     `gorm:"type:varchar(100);not null"`
	Description string
	Password    string
	Handle      string `gorm:"type:varchar(20);not null"`
	Active      bool
	IsPrivate   bool
	OwnerID     uint
	OwnerType   string
}

// HashPassword : hashing the password
func (iban *Iban) HashPassword() {
	hash, err := bcrypt.GenerateFromPassword([]byte(iban.Password), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	iban.Password = string(hash)
}

// ComparePassword : compare the password
func (iban *Iban) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(iban.Password), []byte(password))

	return err == nil
}

// Check Handle before create or update = must be add as index to db
func (iban *Iban) CheckHandle(tx *gorm.DB) (exist bool) {
	var ibans []Iban
	tx.Where("owner_id = ? AND handle = ?", iban.OwnerID, iban.Handle).Find(&ibans)
	for _, tmp := range ibans {
		if iban.Handle == tmp.Handle && iban.IbanID != tmp.IbanID {
			exist = true
		}
	}
	return
}

// BeforeSave Callback
func (iban *Iban) BeforeSave(tx *gorm.DB) (err error) {
	if iban.CheckHandle(tx) {
		err = fmt.Errorf("handle already exist")
	}
	return
}

func (iban *Iban) Validate(db *gorm.DB) {
	if strings.TrimSpace(iban.Text) == "" {
		db.AddError(fmt.Errorf("you have to provide IBAN"))
	} else if strings.TrimSpace(iban.Handle) == "" {
		db.AddError(fmt.Errorf("you have to provide handle"))
	} else if iban.IsPrivate && strings.TrimSpace(iban.Password) == "" {
		db.AddError(fmt.Errorf("you have to provide password"))
	}
}
