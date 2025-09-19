package resolvers

import (
	"context"
	"fmt"
	"strings"

	"github.com/tapsilat/iban.im/config"
	"github.com/tapsilat/iban.im/handler"
	"github.com/tapsilat/iban.im/model"
)

// IbanNew mutation creates iban
func (r *Resolvers) IbanNew(ctx context.Context, args IbanNewMutationArgs) (*IbanNewResponse, error) {
	args.Handle = strings.ToLower(args.Handle)
	UserID := ctx.Value(handler.ContextKey("UserID"))
	if UserID == nil {
		msg := "Not Authorized"
		return &IbanNewResponse{Status: false, Msg: &msg, Iban: nil}, nil
	}
	// userid,_:= strconv.Atoi(UserID.(string))
	userid, _ := UserID.(int)
	fmt.Printf("UserID: %+v, userid: %d\n", UserID, userid)
	if r.HandleCheck(userid, args.Handle) {
		msg := "Same Handle used : " + args.Handle
		return &IbanNewResponse{Status: false, Msg: &msg, Iban: nil}, nil
	}

	// Basic validations (replacing removed qor/validations callbacks)
	if strings.TrimSpace(args.Text) == "" {
		msg := "you have to provide IBAN"
		return &IbanNewResponse{Status: false, Msg: &msg, Iban: nil}, nil
	}
	if strings.TrimSpace(args.Handle) == "" {
		msg := "you have to provide handle"
		return &IbanNewResponse{Status: false, Msg: &msg, Iban: nil}, nil
	}
	if args.IsPrivate && strings.TrimSpace(args.Password) == "" {
		msg := "you have to provide password"
		return &IbanNewResponse{Status: false, Msg: &msg, Iban: nil}, nil
	}

	IbanNew := model.Iban{Text: args.Text, Password: args.Password, Handle: args.Handle, OwnerID: uint(userid), IsPrivate: args.IsPrivate}
	if args.IsPrivate {
		IbanNew.HashPassword()
	}
	if err := config.DB.Create(&IbanNew).Error; err != nil {
		msg := err.Error()
		return &IbanNewResponse{Status: false, Msg: &msg, Iban: nil}, err
	}

	return &IbanNewResponse{Status: true, Msg: nil, Iban: &IbanResponse{i: &IbanNew}}, nil
}

// checks if this handle used for the user
func (r *Resolvers) HandleCheck(userid int, handle string) bool {
	handleStatus := false
	ibans := r.FindIbanByOwner(userid)
	fmt.Printf("ibans: %+v\n", ibans)
	for _, iban := range ibans {
		fmt.Println(iban.Handle)
		if handle == strings.ToLower(iban.Handle) {
			fmt.Println("Same handle found")
			handleStatus = true
			break
		}
	}
	return handleStatus
}

type IbanNewMutationArgs struct {
	Text        string
	Description string
	Password    string
	Handle      string
	IsPrivate   bool
}

// IbanNewResponse is the response type
type IbanNewResponse struct {
	Status bool
	Msg    *string
	Iban   *IbanResponse
}

// Ok for IbanNewResponse
func (r *IbanNewResponse) Ok() bool {
	return r.Status
}

// Error for IbanNewResponse
func (r *IbanNewResponse) Error() *string {
	return r.Msg
}
