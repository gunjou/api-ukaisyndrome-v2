package cdn

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
)

type Service struct {
	StorageRoot string
	BaseURL     string
}

// =================================================
// GET IMAGES GENERAL (MENTOR & MODUL)
// =================================================
func (s *Service) getImages(folder string) ([]ImageDTO, error) {

	path := filepath.Join(s.StorageRoot, folder)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []ImageDTO

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		url := s.BaseURL + "/" + folder + "/" + file.Name()

		result = append(result, ImageDTO{
			Image: url,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Image < result[j].Image
	})

	return result, nil
}

// =================================================
// MENTOR IMAGES
// =================================================
func (s *Service) GetMentorImages() ([]ImageDTO, error) {
	return s.getImages("mentor")
}

// =================================================
// MODUL IMAGES
// =================================================
func (s *Service) GetModulImages() ([]ImageDTO, error) {
	return s.getImages("modul")
}

// =================================================
// GET ADS
// =================================================
func (s *Service) GetAds() (*AdsDTO, error) {

	adsPath := filepath.Join(s.StorageRoot, "ads", "ads.json")

	file, err := os.ReadFile(adsPath)
	if err != nil {
		return nil, err
	}

	var ads []adsMeta

	err = json.Unmarshal(file, &ads)
	if err != nil {
		return nil, err
	}

	for _, ad := range ads {

		if !ad.Active {
			continue
		}

		return &AdsDTO{
			Image: s.BaseURL + "/ads/" + ad.File,
			Link:  ad.Link,
		}, nil
	}

	return nil, errors.New("no active ads found")
}