package resolvers

import (
	"github.com/tapsilat/iban.im/config"
	"github.com/tapsilat/iban.im/model"
	"gorm.io/gorm"
)

// SignUp mutation creates user
func (r *Resolvers) SignUp(args signUpMutationArgs) (*SignUpResponse, error) {

	newUser := model.User{Email: args.Email, Password: args.Password, FirstName: args.FirstName, LastName: args.LastName, Handle: args.Handle}

	var existing model.User
	err := config.DB.Where("email = ? or handle = ?", args.Email, args.Handle).First(&existing).Error
	if err == nil {
		msg := "Already signed up"
		return &SignUpResponse{Status: false, Msg: &msg, User: nil}, nil
	}
	// if not found, err will be gorm.ErrRecordNotFound; any other error should be returned
	if err != gorm.ErrRecordNotFound {
		msg := "lookup error"
		return &SignUpResponse{Status: false, Msg: &msg, User: nil}, err
	}

	newUser.HashPassword()
	if err := config.DB.Create(&newUser).Error; err != nil {
		msg := "create error"
		return &SignUpResponse{Status: false, Msg: &msg, User: nil}, nil
	}

	return &SignUpResponse{Status: true, Msg: nil, User: &UserResponse{u: &newUser}}, nil
}

type signUpMutationArgs struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Handle    string
	Visible   bool
}

// SignUpResponse is the response type
type SignUpResponse struct {
	Status bool
	Msg    *string
	User   *UserResponse
}

// Ok for SignUpResponse
func (r *SignUpResponse) Ok() bool {
	return r.Status
}

// Error for SignUpResponse
func (r *SignUpResponse) Error() *string {
	return r.Msg
}
