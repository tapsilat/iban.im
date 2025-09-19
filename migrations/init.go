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

	defer d.Close()

	d.DropTableIfExists(&model.User{})
	d.DropTableIfExists(&model.Iban{})
	d.DropTableIfExists(&model.Group{})
	d.CreateTable(&model.User{})
	d.CreateTable(&model.Iban{})
	d.CreateTable(&model.Group{})
}
