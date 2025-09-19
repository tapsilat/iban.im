package resolvers

import (
	"context"
	"fmt"

	"github.com/tapsilat/iban.im/config"

	"github.com/tapsilat/iban.im/handler"
	"github.com/tapsilat/iban.im/model"
	// "github.com/tapsilat/iban.im/tools"
)

// GetMyProfile resolver
func (r *Resolvers) GetMyProfile(ctx context.Context) (*GetMyProfileResponse, error) {
	UserID := ctx.Value(handler.ContextKey("UserID"))
	// tools.GetContextDetails(ctx)
	if UserID == nil {
		msg := "Not Authorized"
		return &GetMyProfileResponse{Status: false, Msg: &msg, User: nil}, fmt.Errorf("%s", msg)
	}

	user := model.User{}
	if err := config.DB.First(&user, UserID).Error; err != nil {
		msg := "Not found"
		return &GetMyProfileResponse{Status: false, Msg: &msg, User: nil}, fmt.Errorf("%s", msg)
	}
	return &GetMyProfileResponse{Status: true, Msg: nil, User: &UserResponse{u: &user}}, nil
}

// GetMyProfileResponse is the response type
type GetMyProfileResponse struct {
	Status bool
	Msg    *string
	User   *UserResponse
}

// Ok for GetMyProfileResponse
func (r *GetMyProfileResponse) Ok() bool {
	return r.Status
}

// Error for GetMyProfileResponse
func (r *GetMyProfileResponse) Error() *string {
	return r.Msg
}
