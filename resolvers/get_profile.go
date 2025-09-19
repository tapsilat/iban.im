package resolvers

import (
	"context"
	"fmt"

	"github.com/tapsilat/iban.im/config"
	"github.com/tapsilat/iban.im/model"
)

func (r *Resolvers) getProfileByUserName(userName string) model.User {
	user := model.User{}
	config.DB.Where("handle = ?", userName).First(&user)
	return user
}

func (r *Resolvers) GetProfile(ctx context.Context, args ProfileQueryArgs) (response *GetProfileResponse, err error) {
	user := r.getProfileByUserName(args.Username)
	var ibans []model.Iban
	response = &GetProfileResponse{}
	defer func() {
		if err != nil {
			msg := err.Error()
			response.Msg = &msg
		} else {
			response.Status = true
			response.User = &UserResponse{u: &user}
			var ibansResponse []*IbanResponse
			for _, iban := range ibans {
				tmp := iban
				ibansResponse = append(ibansResponse, &IbanResponse{i: &tmp})
			}
			response.Iban = &ibansResponse
		}
	}()

	if user.UserID == 0 {
		err = fmt.Errorf("user is not exist")
		return
	}
	ibans = r.FindIbanByOwner(int(user.UserID))
	if len(ibans) == 0 {
		err = fmt.Errorf("iban is not exist")
	}
	return
}

// Response for /:username/:alias
type GetProfileResponse struct {
	Status bool
	Msg    *string
	User   *UserResponse
	Iban   *[]*IbanResponse
}

// Ok for GetProfileResponse
func (r *GetProfileResponse) Ok() bool {
	return r.Status
}

// Error for GetProfileResponse
func (r *GetProfileResponse) Error() *string {
	return r.Msg
}

type ProfileQueryArgs struct {
	Username string
}
