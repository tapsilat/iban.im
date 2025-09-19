package resolvers

import (
	"context"

	"github.com/tapsilat/iban.im/handler"
)

// GetMyIbansresolver
func (r *Resolvers) GetMyIbans(ctx context.Context) (*GetMyIbansResponse, error) {
	UserID := ctx.Value(handler.ContextKey("UserID"))
	// tools.GetContextDetails(ctx)
	if UserID == nil {
		msg := "Not Authorized"
		return &GetMyIbansResponse{Status: false, Msg: &msg, Iban: nil}, nil
	}
	userid, _ := UserID.(int)
	ibans := r.FindIbanByOwner(userid)
	var IbansResponse []*IbanResponse
	for _, iban := range ibans {
		tmp := iban
		IbansResponse = append(IbansResponse, &IbanResponse{i: &tmp})
	}

	return &GetMyIbansResponse{Status: true, Msg: nil, Iban: &IbansResponse}, nil
}

// GetMyIbansResponse is the response type
type GetMyIbansResponse struct {
	Status bool
	Msg    *string
	Iban   *[]*IbanResponse
}

// Ok for GetMyIbansResponse
func (r *GetMyIbansResponse) Ok() bool {
	return r.Status
}

// Error for GetMyIbansResponse
func (r *GetMyIbansResponse) Error() *string {
	return r.Msg
}
