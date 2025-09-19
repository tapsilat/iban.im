package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	// gorm postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User : Model with injected fields `ID`, `CreatedAt`, `UpdatedAt`
type User struct {
	UserID    uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Email     string     `gorm:"type:varchar(100);not null"`
	Password  string     `gorm:"not null"`
	Handle    string     `gorm:"not null;unique"`
	FirstName string     `gorm:"type:varchar(50);not null"`
	LastName  string     `gorm:"type:varchar(50);not null"`
	Bio       string     `gorm:"type:text"`
	Visible   bool       // visible email address
	Avatar    string
	Verified  bool
	Active    bool
	Admin     bool
	Ibans     []*Iban `gorm:"polymorphic:Owner;"`
}

// HashPassword : hashing the password
func (user *User) HashPassword() {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return
	}

	user.Password = string(hash)
}

// ComparePassword : compare the password
func (user *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	return err == nil

}
