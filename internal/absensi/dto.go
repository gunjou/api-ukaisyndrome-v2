package absensi

import "time"

type AbsensiPesertaDTO struct {
	IDAbsensiPeserta int      `json:"id_absensi_peserta"`
	IDJadwal         int      `json:"id_jadwal"`
	IDPaketKelas     int      `json:"id_paketkelas"`
	NamaKelas        string   `json:"nama_kelas"`
	IDPeserta        int      `json:"id_peserta"`
	NamaPeserta      string   `json:"nama_peserta"`
	NicknamePeserta  *string  `json:"nickname_peserta"`
	Topik            *string  `json:"topik"`
	Catatan          *string  `json:"catatan"`
	StatusKehadiran  string   `json:"status_kehadiran"`
	CheckInAt        string   `json:"check_in_at"`
	Latitude         *float64 `json:"latitude"`
	Longitude        *float64 `json:"longitude"`
	LocationAccuracy *float64 `json:"location_accuracy"`
}

type AbsensiPesertaStatusDTO struct {
	IDAbsensiPeserta *int    `json:"id_absensi_peserta"`
	IDJadwal         int     `json:"id_jadwal"`
	StatusKehadiran  *string `json:"status_kehadiran"`
	CheckInAt        *string `json:"check_in_at"`
	SudahAbsen       bool    `json:"sudah_absen"`
}

type CheckInPesertaRequest struct {
	IDJadwal         int      `json:"id_jadwal"`
	Latitude         *float64 `json:"latitude"`
	Longitude        *float64 `json:"longitude"`
	LocationAccuracy *float64 `json:"location_accuracy"`
}

type MentorCheckInDTO struct {
	IDAbsensiMentor int
	IDJadwal        int
	TypePertemuan   string
	CheckInAt       time.Time
	CheckOutAt      *time.Time
	Latitude        *float64
	Longitude       *float64
	Accuracy        *float64
}