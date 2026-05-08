package tryout

import (
	"api-ukaisyndrome-v2/pkg/timeutil"
	"context"
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}


// =================================================
// GET TRYOUT BY USER
// =================================================
func (r *Repository) GetTryoutByUser(ctx context.Context, userID int) ([]TryoutDTO, error) {

	query := `
		SELECT 
			t.id_tryout,
			t.judul,
			t.jumlah_soal,
			t.durasi,
			t.max_attempt,
			t.access_start_at,
			t.access_end_at
		FROM tryout t
		JOIN to_paketkelas tp ON tp.id_tryout = t.id_tryout
		JOIN paketkelas pk ON pk.id_paketkelas = tp.id_paketkelas
		JOIN pesertakelas p ON p.id_paketkelas = pk.id_paketkelas
		WHERE 
			p.id_user = $1
			AND t.status = 1
			AND tp.status = 1
			AND t.visibility = 'open'
		ORDER BY t.created_at DESC
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TryoutDTO

	for rows.Next() {
		var t TryoutDTO
		var start, end *time.Time

		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.TotalSoal,
			&t.Duration,
			&t.MaxAttempt,
			&start,
			&end,
		)
		if err != nil {
			return nil, err
		}

		t.AccessStartAt = start
		t.AccessEndAt = end

		result = append(result, t)
	}

	return result, nil
}


// =================================================
// GET TRYOUT FOR START
// =================================================
func (r *Repository) GetTryoutForStart(ctx context.Context, userID int, tryoutID int) (*TryoutDTO, error) {

	query := `
		SELECT 
			t.id_tryout,
			t.judul,
			t.jumlah_soal,
			t.durasi,
			t.max_attempt,
			t.access_start_at,
			t.access_end_at
		FROM tryout t
		WHERE 
			t.id_tryout = $2
			AND t.status = 1
			AND t.visibility = 'open'
			AND EXISTS (
				SELECT 1
				FROM to_paketkelas tp
				JOIN paketkelas pk ON pk.id_paketkelas = tp.id_paketkelas
				JOIN pesertakelas p ON p.id_paketkelas = pk.id_paketkelas
				WHERE 
					p.id_user = $1
					AND tp.id_tryout = t.id_tryout
					AND tp.status = 1
			)
	`

	row := r.DB.QueryRow(ctx, query, userID, tryoutID)

	var t TryoutDTO
	err := row.Scan(
		&t.ID,
		&t.Title,
		&t.TotalSoal,
		&t.Duration,
		&t.MaxAttempt,
		&t.AccessStartAt,
		&t.AccessEndAt,
	)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

// =================================================
// COUNT ATTEMPT
// =================================================
func (r *Repository) CountAttempt(ctx context.Context, userID int, tryoutID int) (int, error) {

	query := `
		SELECT COUNT(*)
		FROM hasiltryout
		WHERE id_user = $1 AND id_tryout = $2
	`

	var count int
	err := r.DB.QueryRow(ctx, query, userID, tryoutID).Scan(&count)
	return count, err
}

// =================================================
// CREATE ATTEMPT
// =================================================
func (r *Repository) InsertAttempt(ctx context.Context, userID int, tryoutID int, attemptToken string, attemptKe int) error {

	query := `
		INSERT INTO hasiltryout (
			id_tryout,
			id_user,
			attempt_token,
			attempt_ke,
			start_time,
			tanggal_pengerjaan,
			status_pengerjaan,
			jawaban_user,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6, 'ongoing', '{}'::jsonb, 1)
	`

	now := timeutil.Now()

	_, err := r.DB.Exec(ctx, query, tryoutID, userID, attemptToken, attemptKe, now, now)
	return err
}


// =================================================
// GET ATTEMPT BY TOKEN
// =================================================
func (r *Repository) GetAttempt(ctx context.Context, attemptToken string, userID int) (*AttemptEntity, error) {

	query := `
		SELECT 
			h.id_tryout,
			h.start_time,
			t.durasi,
			h.status_pengerjaan
		FROM hasiltryout h
		JOIN tryout t ON t.id_tryout = h.id_tryout
		WHERE 
			h.attempt_token = $1
			AND h.id_user = $2
			AND h.status = 1
	`

	var a AttemptEntity

	err := r.DB.QueryRow(ctx, query, attemptToken, userID).Scan(
		&a.TryoutID,
		&a.StartTime,
		&a.Duration,
		&a.Status,
	)

	if err != nil {
		return nil, err
	}

	return &a, nil
}


// =================================================
// GET SOAL BY TRYOUT
// =================================================
func (r *Repository) GetSoalByTryout(ctx context.Context, tryoutID int) ([]SoalDTO, error) {

	query := `
		SELECT 
			id_soaltryout,
			nomor_urut,
			pertanyaan,
			pilihan_a,
			pilihan_b,
			pilihan_c,
			pilihan_d,
			pilihan_e
		FROM soaltryout
		WHERE 
			id_tryout = $1
			AND status = 1
		ORDER BY nomor_urut ASC
	`

	rows, err := r.DB.Query(ctx, query, tryoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []SoalDTO

	for rows.Next() {
		var s SoalDTO
		var a, b, c, d, e string

		err := rows.Scan(
			&s.ID,
			&s.Nomor,
			&s.Pertanyaan,
			&a, &b, &c, &d, &e,
		)
		if err != nil {
			return nil, err
		}

		s.Pilihan = map[string]string{
			"A": a,
			"B": b,
			"C": c,
			"D": d,
			"E": e,
		}

		result = append(result, s)
	}

	return result, nil
}


// =================================================
// SAVE ANSWER (FOR AUTOSAVE)
// =================================================
func (r *Repository) SaveAnswers(
	ctx context.Context,
	attemptToken string,
	userID int,
	answers map[string]AnswerPayload,
) error {

	query := `
		UPDATE hasiltryout
		SET 
			jawaban_user = COALESCE(jawaban_user, '{}'::jsonb) || $1::jsonb,
			updated_at = $4
		WHERE 
			attempt_token = $2
			AND id_user = $3
			AND status = 1
	`

	jsonData, err := json.Marshal(answers)
	if err != nil {
		return err
	}

	now := timeutil.Now()

	_, err = r.DB.Exec(ctx, query, jsonData, attemptToken, userID, now)
	return err
}


// =================================================
// GET ANSWER KEY
// =================================================
func (r *Repository) GetAnswerKey(ctx context.Context, tryoutID int) (map[int]string, error) {

	query := `
		SELECT id_soaltryout, jawaban_benar
		FROM soaltryout
		WHERE id_tryout = $1 AND status = 1
	`

	rows, err := r.DB.Query(ctx, query, tryoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]string)

	for rows.Next() {
		var id int
		var answer string

		if err := rows.Scan(&id, &answer); err != nil {
			return nil, err
		}

		result[id] = answer
	}

	return result, nil
}


// =================================================
// GET USER ANSWERS
// =================================================
func (r *Repository) GetUserAnswers(
	ctx context.Context,
	attemptToken string,
	userID int,
) (map[string]AnswerPayload, error) {

	query := `
		SELECT jawaban_user
		FROM hasiltryout
		WHERE attempt_token = $1 AND id_user = $2
	`

	var jsonData []byte

	err := r.DB.QueryRow(ctx, query, attemptToken, userID).Scan(&jsonData)
	if err != nil {
		return nil, err
	}

	var result map[string]AnswerPayload

	if len(jsonData) > 0 {
		json.Unmarshal(jsonData, &result)
	}

	return result, nil
}


// =================================================
// SUBMIT TRYOUT
// =================================================
func (r *Repository) SubmitResult(
	ctx context.Context,
	attemptToken string,
	userID int,
	benar, salah, kosong, ragu int,
	nilai float64,
) error {

	query := `
		UPDATE hasiltryout
		SET 
			benar = $1,
			salah = $2,
			kosong = $3,
			ragu_ragu = $4,
			nilai = $5,
			status_pengerjaan = 'submitted',
			end_time = $6,
			updated_at = $7
		WHERE attempt_token = $8 AND id_user = $9
	`

	now := timeutil.Now()

	_, err := r.DB.Exec(
		ctx,
		query,
		benar,
		salah,
		kosong,
		ragu,
		nilai,
		now,
		now,
		attemptToken,
		userID,
	)

	return err
}


// =================================================
// GET USER REPORTS
// =================================================
func (r *Repository) GetUserReports(ctx context.Context, userID int) ([]TryoutReportDTO, error) {

	query := `
		SELECT 
			h.attempt_token,
			t.id_tryout,
			t.judul,
			h.attempt_ke,
			h.nilai,
			h.benar,
			h.salah,
			h.kosong,
			h.ragu_ragu,
			h.tanggal_pengerjaan
		FROM hasiltryout h
		JOIN tryout t ON t.id_tryout = h.id_tryout
		WHERE 
			h.id_user = $1
			AND h.status = 1
			AND h.status_pengerjaan = 'submitted'
		ORDER BY h.tanggal_pengerjaan DESC
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TryoutReportDTO

	for rows.Next() {
		var rDTO TryoutReportDTO
		var score sql.NullFloat64

		err := rows.Scan(
			&rDTO.AttemptToken,
			&rDTO.TryoutID,
			&rDTO.Title,
			&rDTO.AttemptKe,
			&score,
			&rDTO.Benar,
			&rDTO.Salah,
			&rDTO.Kosong,
			&rDTO.RaguRagu,
			&rDTO.Tanggal,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, rDTO)
	}

	return result, nil
}



// =================================================
// GET REVIEW
// =================================================
func (r *Repository) GetReview(ctx context.Context, attemptToken string, userID int) ([]ReviewDTO, error) {

	query := `
		SELECT 
			s.id_soaltryout,
			s.nomor_urut,
			s.pertanyaan,
			s.pilihan_a,
			s.pilihan_b,
			s.pilihan_c,
			s.pilihan_d,
			s.pilihan_e,
			s.jawaban_benar,
			s.pembahasan,
			h.jawaban_user
		FROM hasiltryout h
		JOIN soaltryout s ON s.id_tryout = h.id_tryout
		WHERE 
			h.attempt_token = $1
			AND h.id_user = $2
			AND h.status_pengerjaan = 'submitted'
		ORDER BY s.nomor_urut ASC
	`

	rows, err := r.DB.Query(ctx, query, attemptToken, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ReviewDTO

	for rows.Next() {
		var rDTO ReviewDTO
		var a, b, c, d, e string
		var correct string
		var pembahasan string
		var jsonData []byte

		err := rows.Scan(
			&rDTO.ID,
			&rDTO.Nomor,
			&rDTO.Pertanyaan,
			&a, &b, &c, &d, &e,
			&correct,
			&pembahasan,
			&jsonData,
		)
		if err != nil {
			return nil, err
		}

		// decode user answers
		type UserAnswer struct {
			Answer string `json:"answer"`
			Ragu   bool   `json:"ragu"`
		}

		var userAnswers map[string]UserAnswer
		json.Unmarshal(jsonData, &userAnswers)

		// ambil jawaban user
		ua, ok := userAnswers[strconv.Itoa(rDTO.ID)]

		userAns := ""
		isRagu := false

		if ok {
			userAns = ua.Answer
			isRagu = ua.Ragu
		}

		// tentukan status
		status := "kosong"
		if userAns != "" {
			if userAns == correct {
				status = "benar"
			} else {
				status = "salah"
			}
		}

		rDTO.Pilihan = map[string]string{
			"A": a,
			"B": b,
			"C": c,
			"D": d,
			"E": e,
		}

		rDTO.UserAnswer = userAns
		rDTO.CorrectAnswer = correct
		rDTO.Pembahasan = pembahasan
		rDTO.Status = status
		rDTO.IsRagu = isRagu

		result = append(result, rDTO)
	}

	return result, nil
}