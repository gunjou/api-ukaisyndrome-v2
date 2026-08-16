package absensi

import (
	"context"
	"errors"
)

type Service struct {
	Repo *Repository
}

var (
	ErrJadwalNotFound    = errors.New("jadwal tidak ditemukan")
	ErrMentorNotStarted  = errors.New("mentor belum melakukan check-in")
	ErrMeetingFinished   = errors.New("pertemuan sudah selesai")
	ErrCheckInExpired    = errors.New("batas waktu check-in peserta telah berakhir")
	ErrLocationRequired  = errors.New("lokasi wajib dikirim untuk pertemuan offline")
	ErrLocationUnavailable = errors.New("lokasi mentor tidak tersedia")
	ErrOutsideRadius     = errors.New("lokasi peserta berada di luar radius")
	ErrAlreadyCheckedIn  = errors.New("peserta sudah melakukan check-in")
)

func (s *Service) CheckInPeserta(
	ctx context.Context,
	userID int,
	req CheckInPesertaRequest,
) (*AbsensiPesertaDTO, error) {
	exists, err := s.Repo.IsJadwalPeserta(ctx, userID, req.IDJadwal)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrJadwalNotFound
	}

	mentor, err := s.Repo.GetMentorCheckIn(ctx, req.IDJadwal)
	if err != nil {
		return nil, err
	}

	if mentor == nil || mentor.CheckInAt.IsZero() {
		return nil, ErrMentorNotStarted
	}

	if mentor.CheckOutAt != nil {
		return nil, ErrMeetingFinished
	}

	if isCheckInExpired(mentor.CheckInAt) {
		return nil, ErrCheckInExpired
	}

	if mentor.TypePertemuan == "OFFLINE" {
		if req.Latitude == nil || req.Longitude == nil {
			return nil, ErrLocationRequired
		}

		if mentor.Latitude == nil ||
			mentor.Longitude == nil ||
			mentor.Accuracy == nil {
			return nil, ErrLocationUnavailable
		}

		distance := calculateDistance(
			*mentor.Latitude,
			*mentor.Longitude,
			*req.Latitude,
			*req.Longitude,
		)

		if distance > *mentor.Accuracy {
			return nil, ErrOutsideRadius
		}
	}

	alreadyCheckedIn, err := s.Repo.HasAbsensiPeserta(
		ctx,
		userID,
		req.IDJadwal,
	)
	if err != nil {
		return nil, err
	}

	if alreadyCheckedIn {
		return nil, ErrAlreadyCheckedIn
	}

	return s.Repo.InsertCheckIn(ctx, userID, req)
}

func (s *Service) GetAbsensiStatusByJadwal(ctx context.Context, userID int, jadwalID int) (*AbsensiPesertaStatusDTO, error) {
	return s.Repo.GetAbsensiStatusByJadwal(ctx, userID, jadwalID)
}

func (s *Service) GetAbsensiPeserta(ctx context.Context, userID int) ([]AbsensiPesertaDTO, error) {
	return s.Repo.GetAbsensiPeserta(ctx, userID)
}

func (s *Service) GetAbsensiPesertaByID(ctx context.Context, userID, absensiID int) (*AbsensiPesertaDTO, error) {
	return s.Repo.GetAbsensiByID(ctx, userID, absensiID)
}