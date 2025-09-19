package main

import (
	"log"

	"github.com/tapsilat/iban.im/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBOld *gorm.DB
var DBNew *gorm.DB

func init() {
	var err error
	DBOld, err = gorm.Open(mysql.Open("iban_p8nPjkfKO0M:T1os3vVUiQlZeEw@/iban_pushecommerce_com?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DBNew, err = gorm.Open(postgres.Open("host=localhost port=5432 user=postgres dbname=ibanim password=Ahmety61+- sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DBNew.AutoMigrate(&model.User{}, &model.Group{}, &model.Iban{})
}

func main() {
	migrateUser()
	migrateIban()
	migrateGroups()
}

func migrateUser() {
	var users []model.User
	DBOld.Limit(-1).Find(&users)
	for _, user := range users {
		DBNew.Create(&user)
	}
}

func migrateIban() {
	var ibans []model.Iban
	DBOld.Limit(-1).Find(&ibans)
	for _, iban := range ibans {
		DBNew.Create(&iban)
	}
}

func migrateGroups() {
	var groups []model.Group
	DBOld.Limit(-1).Find(&groups)
	for _, group := range groups {
		DBNew.Create(&group)
	}
}
