package module

import (
	"context"
	"errors"
)

type Service struct {
	Repo *Repository
}

func (s *Service) GetModul(ctx context.Context, userID int, role string) ([]ModulDTO, error) {

	if role != "peserta" {
		return nil, errors.New("forbidden")
	}

	return s.Repo.GetModulByUser(ctx, userID)
}