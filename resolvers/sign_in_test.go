package resolvers

import (
	"testing"

	"github.com/tapsilat/iban.im/db"
	"github.com/tapsilat/iban.im/model"
)

func TestSignIn(t *testing.T) {
	db, err := db.ConnectDB("db/database.sqlite")
	if err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	defer db.DB.Close()

	user := model.User{}
	db.DB.Where("email = ?", "notexisting@test.com").First(&user)

	t.Log(user.UserID)
}
