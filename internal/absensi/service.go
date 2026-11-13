package absensi

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"api-ukaisyndrome-v2/pkg/timeutil"
)

type Service struct {
	Repo *Repository
}

var (
	ErrJadwalNotFound          = errors.New("Jadwal tidak ditemukan atau Anda tidak terdaftar pada jadwal ini. Hubungi admin bila seharusnya terdaftar.")
	ErrMentorNotStarted        = errors.New("Mentor belum memulai kelas. Coba lagi setelah mentor melakukan check-in.")
	ErrMeetingFinished         = errors.New("Mentor sudah menutup kelas (check-out). Silakan ajukan absensi susulan.")
	ErrCheckInExpired          = errors.New("Batas waktu check-in sudah lewat. Silakan ajukan absensi susulan setelah jadwal berakhir.")
	ErrLocationRequired        = errors.New("Lokasi wajib dikirim untuk pertemuan offline. Aktifkan lokasi lalu coba lagi.")
	ErrLocationInvalid         = errors.New("Koordinat lokasi tidak valid. Ambil ulang lokasi Anda lalu coba lagi.")
	ErrLocationUnavailable     = errors.New("Lokasi check-in mentor tidak tersedia, sehingga jarak tidak dapat dihitung. Hubungi admin.")
	ErrGPSAccuracyLow          = errors.New("Akurasi GPS terlalu rendah. Pindah ke area terbuka, aktifkan lokasi presisi, lalu coba lagi.")
	ErrOutsideRadius           = errors.New("Anda berada di luar area kelas")
	ErrAlreadyCheckedIn        = errors.New("Anda sudah tercatat hadir pada jadwal ini.")
	ErrSusulanBelumSelesai     = errors.New("Absensi susulan baru dapat diajukan setelah mentor check-out atau jadwal berakhir.")
	ErrSusulanBuktiWajib       = errors.New("Foto bukti kehadiran wajib dilampirkan.")
	ErrSusulanSudahAda         = errors.New("Anda sudah memiliki absensi atau pengajuan susulan pada jadwal ini.")
	ErrSusulanAlasanTidakValid = errors.New("Alasan absensi susulan terlalu panjang (maksimal 2000 karakter).")
)

func isOffline(typePertemuan string) bool {
	return strings.EqualFold(strings.TrimSpace(typePertemuan), "OFFLINE")
}

// CekLokasiCheckIn menjalankan validasi check-in tanpa menyimpan absensi.
// Dipakai juga oleh CheckInPeserta agar indikator di frontend selalu sama dengan hasil check-in.
func (s *Service) CekLokasiCheckIn(
	ctx context.Context,
	userID int,
	req CheckInPesertaRequest,
) error {
	// ONLINE dicek via paket kelas, OFFLINE via kelompok belajar.
	exists, err := s.Repo.IsJadwalPeserta(ctx, userID, req.IDJadwal)
	if err != nil {
		return err
	}
	if !exists {
		return ErrJadwalNotFound
	}

	mentor, err := s.Repo.GetMentorCheckIn(ctx, req.IDJadwal)
	if err != nil {
		return err
	}

	if mentor == nil || mentor.CheckInAt.IsZero() {
		return ErrMentorNotStarted
	}

	if mentor.CheckOutAt != nil {
		return ErrMeetingFinished
	}

	if isCheckInExpired(mentor.CheckInAt) {
		return ErrCheckInExpired
	}

	if isOffline(mentor.TypePertemuan) {
		if req.Latitude == nil || req.Longitude == nil {
			return ErrLocationRequired
		}

		if !isValidLocation(*req.Latitude, *req.Longitude, req.LocationAccuracy) {
			return ErrLocationInvalid
		}

		if accuracyOrDefault(req.LocationAccuracy) > pesertaMaxRadiusMeter {
			return ErrGPSAccuracyLow
		}

		if mentor.Latitude == nil || mentor.Longitude == nil {
			return ErrLocationUnavailable
		}

		distance := calculateDistance(
			*mentor.Latitude,
			*mentor.Longitude,
			*req.Latitude,
			*req.Longitude,
		)

		radius := pesertaRadius(mentor.Accuracy, req.LocationAccuracy)
		if distance > radius {
			return outsideRadiusError(
				*req.Latitude, *req.Longitude,
				*mentor.Latitude, *mentor.Longitude,
				distance, radius,
			)
		}
	}

	return nil
}

func (s *Service) CheckInPeserta(
	ctx context.Context,
	userID int,
	req CheckInPesertaRequest,
) (*AbsensiPesertaDTO, error) {
	if err := s.CekLokasiCheckIn(ctx, userID, req); err != nil {
		return nil, err
	}

	alreadyCheckedIn, err := s.Repo.HasAbsensiPeserta(ctx, userID, req.IDJadwal)
	if err != nil {
		return nil, err
	}
	if alreadyCheckedIn {
		return nil, ErrAlreadyCheckedIn
	}

	data, err := s.Repo.InsertCheckIn(ctx, userID, req)
	if err != nil {
		// Dua request check-in bersamaan: yang kedua ditolak oleh unique index.
		if errors.Is(err, ErrAbsensiDuplikat) {
			return nil, ErrAlreadyCheckedIn
		}
		return nil, err
	}

	return data, nil
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

func (s *Service) AjukanAbsensiSusulan(
	ctx context.Context,
	userID int,
	jadwalID int,
	alasan string,
	filename string,
	contentType string,
	fileData []byte,
) (*AbsensiSusulanDTO, error) {
	if jadwalID <= 0 {
		return nil, ErrJadwalNotFound
	}

	if len(fileData) == 0 {
		return nil, ErrSusulanBuktiWajib
	}

	// Batasi panjang alasan agar input tidak berlebihan.
	if len(alasan) > 2000 {
		return nil, ErrSusulanAlasanTidakValid
	}

	exists, err := s.Repo.IsJadwalPeserta(ctx, userID, jadwalID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrJadwalNotFound
	}

	mentor, err := s.Repo.GetMentorCheckIn(ctx, jadwalID)
	if err != nil {
		return nil, err
	}

	// Mentor yang lupa check-out atau tidak hadir tidak boleh menutup jalur susulan.
	if mentor == nil || mentor.CheckOutAt == nil {
		ended, err := s.Repo.IsJadwalBerakhir(ctx, jadwalID)
		if err != nil {
			return nil, err
		}
		if !ended {
			return nil, ErrSusulanBelumSelesai
		}
	}

	alreadyCheckedIn, err := s.Repo.HasAbsensiPeserta(ctx, userID, jadwalID)
	if err != nil {
		return nil, err
	}
	if alreadyCheckedIn {
		return nil, ErrSusulanSudahAda
	}

	// Upload dilakukan setelah seluruh validasi bisnis awal lolos.
	cdn := NewCDNClient()

	buktiURL, err := cdn.UploadAbsensi(ctx, filename, contentType, fileData)
	if err != nil {
		return nil, fmt.Errorf("upload bukti absensi: %w", err)
	}

	// Gunakan waktu server (timeutil, sama seperti check-in biasa), bukan waktu dari client.
	data, err := s.Repo.InsertAbsensiSusulan(
		ctx,
		userID,
		jadwalID,
		alasan,
		buktiURL,
		timeutil.Now(),
	)
	if err != nil {
		if errors.Is(err, ErrAbsensiDuplikat) {
			return nil, ErrSusulanSudahAda
		}
		return nil, err
	}

	return data, nil
}