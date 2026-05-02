package materi

import (
	"context"
	"errors"
)

type Service struct {
	Repo *Repository
}

func (s *Service) GetMateriPeserta(ctx context.Context, modulID int, materiType *string) ([]MateriDTO, error) {

	// optional: validasi type
	if materiType != nil {
		if *materiType != "video" && *materiType != "document" {
			return nil, errors.New("invalid materi type")
		}
	}

	return s.Repo.GetMateriByModul(ctx, modulID, materiType)
}