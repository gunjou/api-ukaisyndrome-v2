package materi

import (
	"api-ukaisyndrome-v2/pkg/timeutil"
	"context"
	"errors"
)

type Service struct {
	Repo *Repository
}

func (s *Service) GetMateriPeserta(ctx context.Context, userID int, modulID int, materiType *string) ([]MateriDTO, error) {

	// optional: validasi type
	if materiType != nil {
		if *materiType != "video" && *materiType != "document" {
			return nil, errors.New("invalid materi type")
		}
	}

	return s.Repo.GetMateriByModul(ctx, userID, modulID, materiType)
}


func (s *Service) GetMateriPrivatePeserta(ctx context.Context, userID int, materiType *string) ([]MateriPrivateDTO, error) {

	// validasi type
	if materiType != nil {
		if *materiType != "video" && *materiType != "document" {
			return nil, errors.New("invalid materi type")
		}
	}

	return s.Repo.GetMateriPrivateByUser(ctx, userID, materiType)
}


// ============================================================
// MATERI PROGRESS - OPEN / UPSERT
// ============================================================

func (s *Service) OpenMateri(ctx context.Context, userID, materiID int) error {
	if userID <= 0 {
		return errors.New("invalid user id")
	}
	if materiID <= 0 {
		return errors.New("invalid materi id")
	}

	return s.Repo.UpsertMateriProgress(ctx, userID, materiID, timeutil.Now())
}

// ============================================================
// MATERI PROGRESS - RAW DATA
// ============================================================

func (s *Service) GetMateriProgress(ctx context.Context, userID int, modulID *int, page, limit int) ([]MateriProgressDTO, int, error) {
	if userID <= 0 {
		return nil, 0, errors.New("invalid user id")
	}
	if modulID != nil && *modulID <= 0 {
		return nil, 0, errors.New("invalid modul id")
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	return s.Repo.GetMateriProgress(ctx, userID, modulID, page, limit)
}

// ============================================================
// MATERI PROGRESS - MONITORING
// ============================================================

func (s *Service) GetMateriProgressMonitoring(ctx context.Context, userID int, modulID *int) (*MateriProgressMonitoringResponseDTO, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}
	if modulID != nil && *modulID <= 0 {
		return nil, errors.New("invalid modul id")
	}

	summary, modules, err := s.Repo.GetMateriProgressMonitoring(ctx, userID, modulID)
	if err != nil {
		return nil, err
	}

	return &MateriProgressMonitoringResponseDTO{Summary: summary, Modules: modules}, nil
}