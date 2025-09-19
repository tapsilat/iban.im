package resolvers

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go"
	"github.com/tapsilat/iban.im/config"
	"github.com/tapsilat/iban.im/handler"
)

func (r *Resolvers) IbanDelete(ctx context.Context, args IbanDeleteMutationArgs) (response *IbanDeleteResponse, err error) {
	response = &IbanDeleteResponse{}
	iban := r.GetIbanById(args.Id)

	defer func() {
		if err != nil {
			msg := err.Error()
			response.Msg = &msg
		} else {
			response.Status = true
		}
	}()

	userIdStr := ctx.Value(handler.ContextKey("UserID"))
	if userIdStr == nil {
		err = fmt.Errorf("not authorized")
		return
	}

	if iban.OwnerID != uint(userIdStr.(int)) {
		err = fmt.Errorf("not authorized")
		return
	}

	err = config.DB.Delete(&iban).Error

	return
}

type IbanDeleteResponse struct {
	Status bool
	Msg    *string
}

// Ok for IbanDeleteResponse
func (r *IbanDeleteResponse) Ok() bool {
	return r.Status
}

// Error for IbanDeleteResponse
func (r *IbanDeleteResponse) Error() *string {
	return r.Msg
}

// args for delete mutation
type IbanDeleteMutationArgs struct {
	Id graphql.ID
}
