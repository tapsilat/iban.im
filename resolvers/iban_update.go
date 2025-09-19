package resolvers

import (
	"context"
	"strings"

	"github.com/graph-gophers/graphql-go"
	"github.com/tapsilat/iban.im/config"

	"fmt"

	"github.com/tapsilat/iban.im/handler"
	"github.com/tapsilat/iban.im/model"
)

func (r *Resolvers) GetIbanById(id graphql.ID) model.Iban {
	iban := model.Iban{}
	config.DB.Where("iban_id = ?", id).First(&iban)
	return iban
}

// IbanUpdate mutation change profile
func (r *Resolvers) IbanUpdate(ctx context.Context, args IbanUpdateMutationArgs) (response *IbanUpdateResponse, err error) {
	response = &IbanUpdateResponse{}
	iban := r.GetIbanById(args.Id)

	defer func() {
		if err != nil {
			msg := err.Error()
			response.Msg = &msg
		} else {
			response.Status = true
			response.Iban = &IbanResponse{i: &iban}
		}
	}()

	if userID := ctx.Value(handler.ContextKey("UserID")); userID == nil {
		err = fmt.Errorf("not authorized")
		return
	}

	if iban.IbanID == 0 {
		err = fmt.Errorf("iban is not exist")
		return
	}

	iban.Handle = strings.ToLower(args.Handle)
	iban.Text = args.Text
	iban.Description = args.Description

	if args.IsPrivate && args.Password != "" {
		iban.IsPrivate = true
		iban.Password = args.Password
		iban.HashPassword()
	} else if !args.IsPrivate {
		iban.IsPrivate = false
		iban.Password = ""
	}

	err = config.DB.Save(&iban).Error
	return
}

type IbanUpdateMutationArgs struct {
	Text        string
	Description string
	Password    string
	Handle      string
	Id          graphql.ID
	IsPrivate   bool `json:"isPrivate"`
}

// IbanUpdateResponse is the response type
type IbanUpdateResponse struct {
	Status bool
	Msg    *string
	Iban   *IbanResponse
}

// Ok for IbanUpdateResponse
func (r *IbanUpdateResponse) Ok() bool {
	return r.Status
}

// Error for IbanUpdateResponse
func (r *IbanUpdateResponse) Error() *string {
	return r.Msg
}
func (r *Resolvers) FindIbanByHandle(ibans []model.Iban, handle string) model.Iban {
	for _, iban := range ibans {
		fmt.Println(iban.Handle)
		if handle == iban.Handle {
			fmt.Println("Same handle found")
			return iban
		}

	}
	return model.Iban{}
}
func (r *Resolvers) FindIbanByOwner(userID int) []model.Iban {
	var ibans []model.Iban
	// Get all matched records
	config.DB.Where("owner_id = ? AND is_private = false", userID).Find(&ibans)
	return ibans
}
