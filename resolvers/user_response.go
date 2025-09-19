package resolvers

import (
	"strconv"

	graphql "github.com/graph-gophers/graphql-go"

	"fmt"

	"github.com/tapsilat/iban.im/model"
)

// UserResponse is the user response type
type UserResponse struct {
	u *model.User
}

// ID for UserResponse
func (r *UserResponse) ID() graphql.ID {
	id := strconv.Itoa(int(r.u.UserID))
	return graphql.ID(id)
}

// Email for UserResponse
func (r *UserResponse) Email() string {
	return r.u.Email
}

// Password for UserResponse
func (r *UserResponse) Password() string {
	return r.u.Password
}

// FirstName for UserResponse
func (r *UserResponse) FirstName() string {
	return r.u.FirstName
}

// LastName for UserResponse
func (r *UserResponse) LastName() string {
	return r.u.LastName
}

// Visible for UserResponse
func (r *UserResponse) Visible() bool {
	return r.u.Visible
}

// Bio for UserResponse
func (r *UserResponse) Bio() *string {
	return &r.u.Bio
}

// Avatar for UserResponse
func (r *UserResponse) Avatar() *string {
	return &r.u.Avatar
}

// CreatedAt for UserResponse
func (r *UserResponse) CreatedAt() string {
	return r.u.CreatedAt.String()
}

// UpdatedAt for UserResponse
func (r *UserResponse) UpdatedAt() string {
	return r.u.UpdatedAt.String()
}

// Handle for UserResponse
func (r *UserResponse) Handle() string {
	fmt.Printf("ibans: %+v\n", &r.u.Ibans)

	return r.u.Handle
}

// Handle for UserResponse
// func (r *UserResponse) Ibans() *[]*model.Iban {
// 	fmt.Printf("ibans: %+v\n",&r.u.Ibans)
// 	return &r.u.Ibans
// }
