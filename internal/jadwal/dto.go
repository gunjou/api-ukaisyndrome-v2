package jadwal

type JadwalDTO struct {
	IDJadwal int `json:"id_jadwal"`

	// ONLINE  -> akses berdasarkan paket kelas
	// OFFLINE -> akses berdasarkan kelompok belajar
	IDPaketKelas *int    `json:"id_paketkelas"`
	NamaKelas    *string `json:"nama_kelas"`
	IDKelompok   *int    `json:"id_kelompok"`
	NamaKelompok *string `json:"nama_kelompok"`

	IDMentor       int    `json:"id_mentor"`
	NamaMentor     string `json:"nama_mentor"`
	NicknameMentor string `json:"nickname_mentor"`

	Topik   *string `json:"topik"`
	Catatan *string `json:"catatan"`

	Tanggal      string `json:"tanggal"`
	WaktuMulai   string `json:"waktu_mulai"`
	WaktuSelesai string `json:"waktu_selesai"`

	TanggalReschedule      *string `json:"tanggal_reschedule"`
	WaktuMulaiReschedule   *string `json:"waktu_mulai_reschedule"`
	WaktuSelesaiReschedule *string `json:"waktu_selesai_reschedule"`

	TanggalEfektif      string `json:"tanggal_efektif"`
	WaktuMulaiEfektif   string `json:"waktu_mulai_efektif"`
	WaktuSelesaiEfektif string `json:"waktu_selesai_efektif"`

	TypePertemuan string `json:"type_pertemuan"`

	IDAbsensiPeserta  *int    `json:"id_absensi_peserta"`
	StatusKehadiran   *string `json:"status_kehadiran"`
	CheckInAt         *string `json:"check_in_at"`
	JenisAbsensi      *string `json:"jenis_absensi"`
	StatusVerifikasi  *string `json:"status_verifikasi"`
	BuktiKehadiranURL *string `json:"bukti_kehadiran_url"`
	SudahAbsen        bool    `json:"sudah_absen"`
}