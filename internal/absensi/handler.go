package absensi

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"api-ukaisyndrome-v2/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}

const maxAbsensiFileSize = 5 * 1024 * 1024


// CheckInPeserta godoc
// @Summary Check-in peserta
// @Description Peserta melakukan check-in pada jadwal ketika mentor sedang mengajar dan peserta berada dalam radius lokasi mentor
// @Tags Absensi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CheckInPesertaRequest true "Data check-in peserta"
// @Router /absensi/peserta/me/check-in [post]
func (h *Handler) CheckInPeserta(c *fiber.Ctx) error {

	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	var req CheckInPesertaRequest

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, 400, "invalid request body", "BAD_REQUEST", nil)
	}

	if req.IDJadwal <= 0 {
		return response.Error(c, 400, "id_jadwal is required", "BAD_REQUEST", nil)
	}

	data, err := h.Service.CheckInPeserta(
		c.Context(),
		userID,
		req,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrJadwalNotFound):
			return response.Error(c, 404, err.Error(), "NOT_FOUND", nil)

		case errors.Is(err, ErrMentorNotStarted):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)

		case errors.Is(err, ErrMeetingFinished):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)

		case errors.Is(err, ErrCheckInExpired):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)

		case errors.Is(err, ErrLocationRequired):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)

		case errors.Is(err, ErrLocationUnavailable):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)

		case errors.Is(err, ErrOutsideRadius):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)

		case errors.Is(err, ErrAlreadyCheckedIn):
			return response.Error(c, 400, err.Error(), "BAD_REQUEST", nil)
		}

		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}


// GetAbsensiStatusByJadwal godoc
// @Summary Get status absensi peserta berdasarkan jadwal
// @Description Mengecek apakah peserta yang sedang login sudah melakukan check-in pada jadwal tertentu
// @Tags Absensi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_jadwal path int true "ID Jadwal"
// @Router /absensi/peserta/me/jadwal/{id_jadwal} [get]
func (h *Handler) GetAbsensiStatusByJadwal(c *fiber.Ctx) error {
	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	jadwalID, err := c.ParamsInt("id_jadwal")
	if err != nil {
		return response.Error(c, 400, "invalid jadwal id", "BAD_REQUEST", nil)
	}

	data, err := h.Service.GetAbsensiStatusByJadwal(c.Context(), userID, jadwalID)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	if data == nil {
		return response.Error(c, 404, "jadwal tidak ditemukan", "NOT_FOUND", nil)
	}

	return response.Success(c, data)
}


// GetAbsensiPeserta godoc
// @Summary Get absensi peserta
// @Description Get seluruh riwayat absensi peserta yang sedang login
// @Tags Absensi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /absensi/peserta/me [get]
func (h *Handler) GetAbsensiPeserta(c *fiber.Ctx) error {
	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	data, err := h.Service.GetAbsensiPeserta(c.Context(), userID)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	return response.Success(c, data)
}


// GetAbsensiPesertaByID godoc
// @Summary Get detail absensi peserta
// @Description Get detail absensi milik peserta yang sedang login
// @Tags Absensi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id_absensi path int true "ID Absensi Peserta"
// @Router /absensi/peserta/me/{id_absensi} [get]
func (h *Handler) GetAbsensiPesertaByID(c *fiber.Ctx) error {
	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(c, 401, "unauthorized", "UNAUTHORIZED", nil)
	}

	absensiID, err := c.ParamsInt("id_absensi")
	if err != nil {
		return response.Error(c, 400, "invalid absensi id", "BAD_REQUEST", nil)
	}

	data, err := h.Service.GetAbsensiPesertaByID(c.Context(), userID, absensiID)
	if err != nil {
		return response.Error(c, 500, err.Error(), "INTERNAL_ERROR", nil)
	}

	if data == nil {
		return response.Error(c, 404, "absensi tidak ditemukan", "NOT_FOUND", nil)
	}

	return response.Success(c, data)
}



 // AjukanAbsensiSusulan godoc
 // @Summary Ajukan absensi susulan peserta
 // @Description Mengajukan absensi susulan setelah mentor melakukan check-out dengan melampirkan bukti foto kehadiran.
 // @Tags Absensi
 // @Accept mpfd
 // @Produce json
 // @Security BearerAuth
 // @Param id_jadwal formData int true "ID Jadwal Kelas"
 // @Param file formData file true "Foto bukti kehadiran (JPEG, PNG, atau WebP; maksimal 5 MB)"
 // @Param alasan_absensi_susulan formData string false "Alasan mengajukan absensi susulan"
 // @Router /absensi/peserta/me/absensi-susulan [post]
func (h *Handler) AjukanAbsensiSusulan(c *fiber.Ctx) error {
	userID, ok := c.Locals("sub").(int)
	if !ok {
		return response.Error(
			c, 401, "unauthorized", "UNAUTHORIZED", nil,
		)
	}

	jadwalID, err := strconv.Atoi(c.FormValue("id_jadwal"))
	if err != nil || jadwalID <= 0 {
		return response.Error(
			c, 400, "id_jadwal tidak valid", "BAD_REQUEST", nil,
		)
	}

	alasan := strings.TrimSpace(
		c.FormValue("alasan_absensi_susulan"),
	)
	if len(alasan) > 2000 {
		return response.Error(
			c, 400, "alasan maksimal 2000 karakter",
			"BAD_REQUEST", nil,
		)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return response.Error(
			c, 400, "foto bukti kehadiran wajib dilampirkan",
			"BAD_REQUEST", nil,
		)
	}

	if fileHeader.Size <= 0 || fileHeader.Size > maxAbsensiFileSize {
		return response.Error(
			c, 400, "ukuran file harus antara 1 byte dan 5 MB",
			"BAD_REQUEST", nil,
		)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return response.Error(
			c, 400, "file tidak dapat dibaca", "BAD_REQUEST", nil,
		)
	}
	defer file.Close()

	fileData, err := io.ReadAll(
		io.LimitReader(file, maxAbsensiFileSize+1),
	)
	if err != nil {
		return response.Error(
			c, 400, "gagal membaca file", "BAD_REQUEST", nil,
		)
	}

	if len(fileData) == 0 || len(fileData) > maxAbsensiFileSize {
		return response.Error(
			c, 400, "ukuran file tidak valid", "BAD_REQUEST", nil,
		)
	}

	contentType := detectAbsensiImageType(fileData)
	if contentType == "" {
		return response.Error(
			c, 400, "file harus berupa gambar JPEG, PNG, atau WebP",
			"BAD_REQUEST", nil,
		)
	}

	data, err := h.Service.AjukanAbsensiSusulan(
		c.Context(),
		userID,
		jadwalID,
		alasan,
		fileHeader.Filename,
		contentType,
		fileData,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrJadwalNotFound):
			return response.Error(
				c, 404, err.Error(), "NOT_FOUND", nil,
			)

		case errors.Is(err, ErrMentorNotStarted),
			errors.Is(err, ErrSusulanBelumSelesai),
			errors.Is(err, ErrSusulanBuktiWajib),
			errors.Is(err, ErrSusulanSudahAda),
			errors.Is(err, ErrSusulanAlasanTidakValid):
			return response.Error(
				c, 400, err.Error(), "BAD_REQUEST", nil,
			)

		case errors.Is(err, ErrCDNNotConfigured),
			errors.Is(err, ErrCDNUploadFailed):
			return response.Error(
				c, 502, "gagal memproses bukti kehadiran",
				"CDN_ERROR", nil,
			)
		}

		// Detail error internal dicatat di server, bukan dikirim
		// kepada client karena dapat mengandung informasi sensitif.
		return response.Error(
			c, 500, "gagal menyimpan pengajuan absensi",
			"INTERNAL_ERROR", nil,
		)
	}

	return response.Success(c, data)
}

func detectAbsensiImageType(data []byte) string {
	if len(data) < 12 {
		return ""
	}

	// JPEG
	if len(data) >= 3 &&
		data[0] == 0xFF &&
		data[1] == 0xD8 &&
		data[2] == 0xFF {
		return "image/jpeg"
	}

	// PNG
	if len(data) >= 8 &&
		string(data[:8]) == "\x89PNG\r\n\x1a\n" {
		return "image/png"
	}

	// WebP: RIFF....WEBP
	if string(data[:4]) == "RIFF" &&
		string(data[8:12]) == "WEBP" {
		return "image/webp"
	}

	detected := http.DetectContentType(data)
	if detected == "image/jpeg" || detected == "image/png" {
		return detected
	}

	return ""
}
