package resolvers

import (
	"strconv"

	graphql "github.com/graph-gophers/graphql-go"

	"github.com/tapsilat/iban.im/model"
)

// IbanResponse is the user response type
type IbanResponse struct {
	i *model.Iban
}

// UserID for IbanResponse
func (r *IbanResponse) ID() graphql.ID {
	id := strconv.Itoa(int(r.i.IbanID))
	return graphql.ID(id)
}

// UserID for IbanResponse
func (r *IbanResponse) OwnerID() string {
	ownerId := strconv.Itoa(int(r.i.OwnerID))
	return ownerId
}

// Text for IbanResponse
func (r *IbanResponse) Text() string {
	return r.i.Text
}

// Description for IbanResponse
func (r *IbanResponse) Description() *string {
	return &r.i.Description
}

// Password for IbanResponse
func (r *IbanResponse) Password() string {
	return r.i.Password
}

// Handle for IbanResponse
func (r *IbanResponse) Handle() string {
	return r.i.Handle
}

// IsPrivate for IbanResponse
func (r *IbanResponse) IsPrivate() bool {
	return r.i.IsPrivate
}

// CreatedAt for IbanResponse
func (r *IbanResponse) CreatedAt() string {
	return r.i.CreatedAt.String()
}

// UpdatedAt for IbanResponse
func (r *IbanResponse) UpdatedAt() string {
	return r.i.UpdatedAt.String()
}
