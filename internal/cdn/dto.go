package cdn

type ImageDTO struct {
	Image string `json:"image"`
}

type AdsDTO struct {
	Image string `json:"image"`
	Link  string `json:"link"`
}

type adsMeta struct {
	ID     int    `json:"id"`
	File   string `json:"file"`
	Link   string `json:"link"`
	Active bool   `json:"active"`
}