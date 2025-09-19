package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/tapsilat/iban.im/utils"
)

// ContextKey for the userID in context
type ContextKey string

// Authenticate for JWT
func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// var userID *string
		fmt.Println("inside authentication")
		ctx := r.Context()
		userID, err := validateAuthHeader(ctx, r)
		if err != nil {
			// should do something here
			fmt.Println("validate auth header error")
		}

		if userID != nil {
			ctx = context.WithValue(ctx, ContextKey("UserID"), *userID)
		}

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateAuthHeader(ctx context.Context, r *http.Request) (*string, error) {
	_ = ctx //Added to show temporarily used.
	fmt.Println("inside validate auth header")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, nil
	}

	userID, err := utils.ValidateJWT(&tokenString)
	return userID, err
}
