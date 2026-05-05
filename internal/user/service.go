package user

import (
	"api-ukaisyndrome-v2/internal/shared"
	"context"
	"errors"
)

type Service struct {
	Repo *Repository
}


// ==========================
// GetMe returns current logged in user profile based on JWT
// ==========================
func (s *Service) GetMe(ctx context.Context, userID int, role string) (*MeResponse, error) {

	user, err := s.Repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := &MeResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	switch role {

		case "admin":
			// hanya basic info
			return res, nil

		case "mentor":
			res.Nickname = user.Nickname
			return res, nil

		case "peserta":
			res.Nickname = user.Nickname // kalau ada, optional

			classes, err := s.Repo.GetUserClasses(ctx, userID)
			if err != nil {
				return nil, err
			}

			res.Classes = classes
			return res, nil

		default:
			return nil, errors.New("invalid role")
	}
}


// ==========================
// CHANGE PASSWORD
// ==========================
func (s *Service) ChangePassword(ctx context.Context, userID int, req ChangePasswordRequest) error {

	user, err := s.Repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// VALIDASI PASSWORD LAMA
	if !shared.VerifyHash(user.Password, req.OldPassword) {
		return errors.New("old password is incorrect")
	}

	// HASH PASSWORD BARU
	hashed, err := shared.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// UPDATE PASSWORD
	return s.Repo.UpdatePassword(ctx, userID, hashed)
}