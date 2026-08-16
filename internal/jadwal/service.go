package jadwal

import "context"

type Service struct {
	Repo *Repository
}

func (s *Service) GetJadwalPeserta(ctx context.Context, userID int) ([]JadwalDTO, error) {

	return s.Repo.GetJadwalPeserta(
		ctx,
		userID,
	)
}

func (s *Service) GetJadwalPesertaByID(ctx context.Context, userID int, jadwalID int) (*JadwalDTO, error) {

	return s.Repo.GetJadwalPesertaByID(
		ctx,
		userID,
		jadwalID,
	)
}