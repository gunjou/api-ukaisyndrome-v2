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