package jadwal

type JadwalDTO struct {
	IDJadwal              int    `json:"id_jadwal"`
	IDPaketKelas          int    `json:"id_paketkelas"`
	NamaKelas             string `json:"nama_kelas"`
	IDMentor               int    `json:"id_mentor"`
	NamaMentor             string `json:"nama_mentor"`
	NicknameMentor         string `json:"nickname_mentor"`
	Tanggal                string `json:"tanggal"`
	WaktuMulai             string `json:"waktu_mulai"`
	WaktuSelesai           string `json:"waktu_selesai"`
	TanggalReschedule      *string `json:"tanggal_reschedule"`
	WaktuMulaiReschedule   *string `json:"waktu_mulai_reschedule"`
	WaktuSelesaiReschedule *string `json:"waktu_selesai_reschedule"`
	TanggalEfektif         string `json:"tanggal_efektif"`
	WaktuMulaiEfektif      string `json:"waktu_mulai_efektif"`
	WaktuSelesaiEfektif    string `json:"waktu_selesai_efektif"`
	TypePertemuan          string `json:"type_pertemuan"`
}