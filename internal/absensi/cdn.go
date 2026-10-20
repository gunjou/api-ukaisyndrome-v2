package absensi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"time"
)

var (
	ErrCDNNotConfigured = errors.New("konfigurasi CDN belum lengkap")
	ErrCDNUploadFailed  = errors.New("gagal mengunggah bukti ke CDN")
)

type CDNClient struct {
	UploadURL string
	APIKey    string
	HTTP      *http.Client
}

type CDNUploadResponse struct {
	Status bool   `json:"status"`
	URL    string `json:"url"`
	Size   int64  `json:"size"`
	File   string `json:"file"`
}

func NewCDNClient() *CDNClient {
	return &CDNClient{
		UploadURL: strings.TrimRight(
			os.Getenv("CDN_UPLOAD_URL"), "/",
		),
		APIKey: os.Getenv("CDN_API_KEY"),
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *CDNClient) UploadAbsensi(
	ctx context.Context,
	filename string,
	contentType string,
	fileData []byte,
) (string, error) {
	if c.UploadURL == "" || c.APIKey == "" {
		return "", ErrCDNNotConfigured
	}

	if len(fileData) == 0 {
		return "", errors.New("file kosong")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	header := make(textproto.MIMEHeader)
	header.Set(
		"Content-Disposition",
		fmt.Sprintf(
			`form-data; name="file"; filename="%s"`,
			safeFilename(filename),
		),
	)
	header.Set("Content-Type", contentType)

	part, err := writer.CreatePart(header)
	if err != nil {
		return "", fmt.Errorf("buat multipart: %w", err)
	}

	if _, err := part.Write(fileData); err != nil {
		return "", fmt.Errorf("tulis file multipart: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("tutup multipart: %w", err)
	}

	uploadURL := c.UploadURL + "/absensi"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		uploadURL,
		&body,
	)
	if err != nil {
		return "", fmt.Errorf("buat request CDN: %w", err)
	}

	req.Header.Set("X-API-KEY", c.APIKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrCDNUploadFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Jangan kirim body respons CDN ke client.
		_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
		return "", fmt.Errorf(
			"%w: HTTP %d",
			ErrCDNUploadFailed,
			resp.StatusCode,
		)
	}

	var result CDNUploadResponse

	if err := json.NewDecoder(
		io.LimitReader(resp.Body, 64*1024),
	).Decode(&result); err != nil {
		return "", fmt.Errorf(
			"%w: respons tidak valid",
			ErrCDNUploadFailed,
		)
	}

	if !result.Status || result.URL == "" {
		return "", fmt.Errorf(
			"%w: URL file tidak tersedia",
			ErrCDNUploadFailed,
		)
	}

	// Pastikan CDN mengembalikan URL HTTPS dari host yang dipercaya.
	// Jangan menyimpan URL arbitrer dari respons eksternal.
	if !isTrustedCDNURL(result.URL) {
		return "", fmt.Errorf(
			"%w: URL hasil upload tidak valid",
			ErrCDNUploadFailed,
		)
	}

	return result.URL, nil
}

func safeFilename(filename string) string {
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, `"`, "_")
	filename = strings.ReplaceAll(filename, "\r", "")
	filename = strings.ReplaceAll(filename, "\n", "")

	if filename == "" {
		return "bukti-kehadiran"
	}

	return filename
}

func isTrustedCDNURL(rawURL string) bool {
	// Host disesuaikan dengan CDN milik aplikasi ini.
	// Implementasi validasi URL lengkap ada di bawah.
	return strings.HasPrefix(
		rawURL,
		"https://mgf.ukaisyndrome.id/images/ukaisyndrome/",
	)
}
