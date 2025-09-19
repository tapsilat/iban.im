package resolvers

import (
	// "strconv"

	"fmt"

	"github.com/tapsilat/iban.im/config"
	"github.com/tapsilat/iban.im/model"
	"github.com/tapsilat/iban.im/utils"
	// "fmt"
)

// SignIn mutation creates user
func (r *Resolvers) SignIn(args signInMutationArgs) (response *SignInResponse, err error) {
	response = &SignInResponse{}
	user := r.GetUserByEmail(args.Email)
	var tokenString *string
	defer func() {
		if err != nil {
			msg := err.Error()
			response.Msg = &msg
		} else {
			response.Status = true
			response.Token = tokenString
		}
	}()

	if user.UserID == 0 {
		err = fmt.Errorf("not sign up yet")
		return
	}

	if !user.ComparePassword(args.Password) {
		err = fmt.Errorf("password is not correct")
		return
	}

	tokenString, err = utils.SignJWT(&args.Email, &args.Password)
	return
}

func (r *Resolvers) GetUserByEmail(email string) model.User {
	user := model.User{}
	config.DB.Where("email = ?", email).First(&user)
	return user
}

type signInMutationArgs struct {
	Email    string
	Password string
}

// SignInResponse is the response type
type SignInResponse struct {
	Status bool
	Msg    *string
	Token  *string
}

// Ok for SignUpResponse
func (r *SignInResponse) Ok() bool {
	return r.Status
}

// Error for SignUpResponse
func (r *SignInResponse) Error() *string {
	return r.Msg
}
