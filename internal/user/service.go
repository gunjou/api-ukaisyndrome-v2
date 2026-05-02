package user

import (
	"context"
	"errors"
)

type Service struct {
	Repo *Repository
}

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