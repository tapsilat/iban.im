package resolvers

import (
	"context"
	"log"

	"github.com/tapsilat/iban.im/config"
	"github.com/tapsilat/iban.im/handler"
	"github.com/tapsilat/iban.im/model"
)

// DeleteProfile mutation to delete user profile
func (r *Resolvers) DeleteProfile(ctx context.Context, args deleteProfileMutationArgs) (*DeleteProfileResponse, error) {
	userID := ctx.Value(handler.ContextKey("UserID"))

	if userID == nil {
		msg := "Not Authorized"
		return &DeleteProfileResponse{Status: false, Msg: &msg, Message: nil}, nil
	}

	user := model.User{}
	if err := config.DB.First(&user, userID).Error; err != nil {
		msg := "User not found"
		return &DeleteProfileResponse{Status: false, Msg: &msg, Message: nil}, nil
	}

	// Verify password as confirmation
	if !user.ComparePassword(args.ConfirmPassword) {
		msg := "Invalid password confirmation"
		return &DeleteProfileResponse{Status: false, Msg: &msg, Message: nil}, nil
	}

	// Delete all associated IBANs (soft delete)
	if err := config.DB.Where("owner_id = ? AND owner_type = ?", user.UserID, "User").Delete(&model.Iban{}).Error; err != nil {
		msg := "Failed to delete user IBANs"
		log.Printf("Error deleting IBANs for user %d: %v", user.UserID, err)
		return &DeleteProfileResponse{Status: false, Msg: &msg, Message: nil}, err
	}

	// Delete user (soft delete using GORM's DeletedAt)
	if err := config.DB.Delete(&user).Error; err != nil {
		msg := "Failed to delete user profile"
		log.Printf("Error deleting user %d: %v", user.UserID, err)
		return &DeleteProfileResponse{Status: false, Msg: &msg, Message: nil}, err
	}

	// Log the deletion for auditing
	log.Printf("User profile deleted: UserID=%d, Email=%s, Handle=%s", user.UserID, user.Email, user.Handle)

	successMsg := "Profile deleted successfully. Your account and all associated data have been removed."
	return &DeleteProfileResponse{Status: true, Msg: nil, Message: &successMsg}, nil
}

type deleteProfileMutationArgs struct {
	ConfirmPassword string
}

// DeleteProfileResponse is the response type
type DeleteProfileResponse struct {
	Status  bool
	Msg     *string
	Message *string
}

// Ok for DeleteProfileResponse
func (r *DeleteProfileResponse) Ok() bool {
	return r.Status
}

// Error for DeleteProfileResponse
func (r *DeleteProfileResponse) Error() *string {
	return r.Msg
}
