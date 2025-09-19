package main

import (
	"flag"

	"github.com/tapsilat/iban.im/db"
	"github.com/tapsilat/iban.im/model"
)

var env string

func main() {

	flag.StringVar(&env, "env", "localhost", "[localhost docker gitpod]")
	flag.Parse()

	d, err := db.ConnectDB(env)
	if err != nil {
		panic(err)
	}

	sqlDB, err := d.DB.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	// Drop and create tables using Migrator (GORM v2)
	d.Migrator().DropTable(&model.User{}, &model.Iban{}, &model.Group{})
	d.AutoMigrate(&model.User{}, &model.Iban{}, &model.Group{})
}
