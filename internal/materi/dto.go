package materi

type MateriDTO struct {
	ID             int    `json:"id"`
	IDModul        int    `json:"id_modul"`
	Type           string `json:"type"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	IsDownloadable int    `json:"is_downloadable"`
}

type MateriPrivateDTO struct {
	ID             int    `json:"id"`
	Type           string `json:"type"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	IsDownloadable int    `json:"is_downloadable"`
}

// ============================================================
// MATERI PROGRESS
// ============================================================

type MateriProgressDTO struct {
	IDMateriProgress int    `json:"id_materi_progress"`
	IDMateri         int    `json:"id_materi"`
	FirstOpenedAt    string `json:"first_opened_at"`
	LastOpenedAt     string `json:"last_opened_at"`
}

// ============================================================
// MATERI PROGRESS MONITORING
// ============================================================

type MateriProgressMonitoringDTO struct {
	IDModul            int     `json:"id_modul"`
	NamaModul          string  `json:"nama_modul"`
	TotalMateri        int     `json:"total_materi"`
	MateriDibuka       int     `json:"materi_dibuka"`
	MateriBelumDibuka int     `json:"materi_belum_dibuka"`
	ProgressPercentage float64 `json:"progress_percentage"`
}

type MateriProgressMonitoringSummaryDTO struct {
	TotalMateri        int     `json:"total_materi"`
	MateriDibuka       int     `json:"materi_dibuka"`
	MateriBelumDibuka int     `json:"materi_belum_dibuka"`
	ProgressPercentage float64 `json:"progress_percentage"`
}

type MateriProgressMonitoringResponseDTO struct {
	Summary MateriProgressMonitoringSummaryDTO `json:"summary"`
	Modules []MateriProgressMonitoringDTO      `json:"modules"`
}

// ============================================================
// PAGINATION
// ============================================================

type PaginationDTO struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type MateriProgressListResponseDTO struct {
	Data       []MateriProgressDTO `json:"data"`
	Pagination PaginationDTO       `json:"pagination"`
}