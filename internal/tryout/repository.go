package tryout

import (
	"api-ukaisyndrome-v2/internal/shared/score"
	"api-ukaisyndrome-v2/pkg/timeutil"
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}


/* ========================================================================== */
/*                             //SECTION - TRYOUT                             */
/* ========================================================================== */

//ANCHOR - GET TRYOUT BY USER
func (r *Repository) GetTryoutByUser(ctx context.Context, userID int) ([]TryoutDTO, error) {

	query := `
		SELECT 
			t.id_tryout,
			t.judul,
			t.jumlah_soal,
			t.durasi,
			t.max_attempt,

			CASE 
				WHEN t.max_attempt > 0 
				THEN t.max_attempt - COUNT(h.id_hasiltryout)
				ELSE 0
			END as remaining_attempt,

			t.access_start_at,
			t.access_end_at

		FROM tryout t
		JOIN to_paketkelas tp ON tp.id_tryout = t.id_tryout
		JOIN paketkelas pk ON pk.id_paketkelas = tp.id_paketkelas
		JOIN pesertakelas p ON p.id_paketkelas = pk.id_paketkelas

		LEFT JOIN hasiltryout h 
			ON h.id_tryout = t.id_tryout 
			AND h.id_user = p.id_user
			AND h.status = 1

		WHERE 
			p.id_user = $1
			AND t.status = 1
			AND tp.status = 1
			AND t.visibility = 'open'

		GROUP BY 
			t.id_tryout

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
		var remaining int

		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.TotalSoal,
			&t.Duration,
			&t.MaxAttempt,
			&remaining,
			&start,
			&end,
		)
		if err != nil {
			return nil, err
		}

		t.RemainingAttempt = remaining
		t.AccessStartAt = start
		t.AccessEndAt = end

		result = append(result, t)
	}

	return result, nil
}



//ANCHOR - GET TRYOUT FOR START
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

//ANCHOR - COUNT ATTEMPT
func (r *Repository) CountAttempt(ctx context.Context, userID int, tryoutID int) (int, error) {

	query := `
		SELECT COUNT(*)
		FROM hasiltryout
		WHERE id_user = $1 AND id_tryout = $2 AND status = 1
	`

	var count int
	err := r.DB.QueryRow(ctx, query, userID, tryoutID).Scan(&count)
	return count, err
}

//ANCHOR - CREATE ATTEMPT
func (r *Repository) InsertAttempt(ctx context.Context, userID int, tryoutID int, attemptToken string, attemptKe int, duration int,
) error {

	query := `
		INSERT INTO hasiltryout (
			id_tryout, id_user, attempt_token, attempt_ke, start_time, end_time, 
			tanggal_pengerjaan, status_pengerjaan, jawaban_user, status
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			'ongoing',
			'{}'::jsonb,
			1

		)
	`

	now := timeutil.Now()

	endTime := now.Add(time.Duration(duration) * time.Minute)

	_, err := r.DB.Exec(
		ctx,
		query,
		tryoutID,
		userID,
		attemptToken,
		attemptKe,
		now,
		endTime,
		now,
	)

	return err }


//ANCHOR - GET ONGOING ATTEMPTS
func (r *Repository) GetOngoingAttempts(
	ctx context.Context,
	userID int,
	tryoutID int,
) ([]OngoingAttemptEntity, error) {

	query := `
		SELECT

			id_hasiltryout,
			attempt_token,
			attempt_ke,
			start_time,
			end_time

		FROM hasiltryout

		WHERE

			id_user = $1
			AND id_tryout = $2
			AND status = 1
			AND status_pengerjaan = 'ongoing'

		ORDER BY attempt_ke ASC
	`

	rows, err := r.DB.Query(ctx, query, userID, tryoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []OngoingAttemptEntity

	for rows.Next() {

		var item OngoingAttemptEntity

		if err := rows.Scan(
			&item.ID,
			&item.AttemptToken,
			&item.AttemptKe,
			&item.StartTime,
			&item.EndTime,
		); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}



//ANCHOR - GET ATTEMPT BY TOKEN
func (r *Repository) GetAttempt(ctx context.Context, attemptToken string, userID int) (*AttemptEntity, error) {

	query := `
		SELECT 
			h.id_tryout,
			h.start_time,
			h.end_time,
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
		&a.EndTime,
		&a.Duration,
		&a.Status,
	)

	if err != nil {
		return nil, err
	}

	return &a, nil
}



//ANCHOR - GET SOAL BY TRYOUT
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

//ANCHOR - SAVE ANSWERS
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


//ANCHOR - GET ANSWER KEY
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


//ANCHOR - GET USER ANSWERS
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


//ANCHOR - SUBMIT TRYOUT
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


//ANCHOR - UPDATE SUBMITTED ATTEMPT
func (r *Repository) UpdateSubmittedAttempt(
	ctx context.Context,
	id int,
	nilai float64,
	benar int,
	salah int,
	kosong int,
	ragu int,
) error {

	query := `
		UPDATE hasiltryout
		SET

			nilai = $1,
			benar = $2,
			salah = $3,
			kosong = $4,
			ragu_ragu = $5,
			status_pengerjaan = 'submitted',
			end_time = $6,
			updated_at = $7

		WHERE

			id_hasiltryout = $8
	`

	now := timeutil.Now()

	_, err := r.DB.Exec(
		ctx,
		query,
		nilai,
		benar,
		salah,
		kosong,
		ragu,
		now,
		now,
		id,
	)

	return err
}


//ANCHOR - GET ONGOING TRYOUT
func (r *Repository) GetOngoingTryout(
	ctx context.Context,
	userID int,
) ([]OngoingTryoutEntity, error) {

	query := `
		SELECT
			h.id_hasiltryout,
			h.id_tryout,
			t.judul,
			h.attempt_token,
			h.start_time,
			h.end_time,
			h.jawaban_user,
			h.status_pengerjaan
		FROM hasiltryout h
		JOIN tryout t
			ON t.id_tryout = h.id_tryout
		WHERE
			h.id_user = $1
			AND h.status = 1
			AND h.status_pengerjaan = 'ongoing'
		ORDER BY h.start_time DESC
	`

	rows, err := r.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []OngoingTryoutEntity

	for rows.Next() {

		var item OngoingTryoutEntity

		err := rows.Scan(
			&item.IDHasilTryout,
			&item.IDTryout,
			&item.Title,
			&item.AttemptToken,
			&item.StartTime,
			&item.EndTime,
			&item.JawabanUser,
			&item.StatusPengerjaan,
		)
		if err != nil {
			return nil, err
		}

		// NORMALIZE TIME
		item.StartTime = timeutil.NormalizeDBTime(item.StartTime)
		if item.EndTime != nil {
			t := timeutil.NormalizeDBTime(*item.EndTime)
			item.EndTime = &t
		}

		result = append(result, item)
	}

	return result, nil
}
/* ========================== //!SECTION - TRYOUT ========================== */



/* ========================================================================== */
/*                             //SECTION - REPORT                             */
/* ========================================================================== */

//ANCHOR - GET USER REPORTS
func (r *Repository) GetUserReports(ctx context.Context, userID int) ([]TryoutReportDTO, error) {

	query := `
		SELECT 
			h.attempt_token,
			t.id_tryout,
			t.judul,
			h.attempt_ke,
			CASE 
				WHEN t.jumlah_soal > 0 
				THEN (h.benar::float / t.jumlah_soal) * 100
				ELSE 0
			END as score,
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

		err := rows.Scan(
			&rDTO.AttemptToken,
			&rDTO.TryoutID,
			&rDTO.Title,
			&rDTO.AttemptKe,
			&rDTO.Score,
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


//ANCHOR - GET REVIEW
func (r *Repository) GetReview(
	ctx context.Context,
	attemptToken string,
	userID int,
) (string, []ReviewDTO, error) {

	query := `
		SELECT
			t.id_tryout,
			t.judul,
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
		JOIN tryout t
			ON t.id_tryout = h.id_tryout
		JOIN soaltryout s
			ON s.id_tryout = h.id_tryout
		WHERE
			h.attempt_token = $1
			AND h.id_user = $2
			AND h.status_pengerjaan = 'submitted'
			AND s.status = 1
		ORDER BY s.nomor_urut ASC
	`

	rows, err := r.DB.Query(ctx, query, attemptToken, userID)
	if err != nil {
		return "", nil, err
	}
	defer rows.Close()

	var (
		title  string
		result []ReviewDTO
	)

	for rows.Next() {

		var (
			idTryout int
			rowTitle string
		)

		var rDTO ReviewDTO
		var a, b, c, d, e string
		var correct string
		var pembahasan string
		var jsonData []byte

		err := rows.Scan(
			&idTryout,
			&rowTitle,
			&rDTO.ID,
			&rDTO.Nomor,
			&rDTO.Pertanyaan,
			&a,
			&b,
			&c,
			&d,
			&e,
			&correct,
			&pembahasan,
			&jsonData,
		)
		if err != nil {
			return "", nil, err
		}

		// simpan title sekali saja
		if title == "" {
			title = rowTitle
		}

		// decode user answers
		type UserAnswer struct {
			Answer string `json:"answer"`
			Ragu   bool   `json:"ragu"`
		}

		var userAnswers map[string]UserAnswer
		_ = json.Unmarshal(jsonData, &userAnswers)

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

	if err := rows.Err(); err != nil {
		return "", nil, err
	}

	return title, result, nil
}
/* =========================== //!SECTION - REPORT ========================== */


/* ========================================================================== */
/*                //SECTION - ENDPOINT LIST LEADERBOARD SECTION               */
/* ========================================================================== */

//ANCHOR - GET GLOBAL LEADERBOARD
func (r *Repository) GetGlobalLeaderboard(ctx context.Context, tryoutID int) ([]LeaderboardItem, error) {

	query := `
		WITH ranked_attempt AS (

			SELECT
				h.id_user, u.nama, pk.nama_kelas,
		 		h.attempt_ke, h.nilai, h.benar, h.start_time, h.end_time,

				EXTRACT(EPOCH FROM (h.end_time - h.start_time))::BIGINT AS duration,

				ROW_NUMBER() OVER (
					PARTITION BY h.id_user
					ORDER BY
						h.nilai DESC,
						h.benar DESC,
						(h.end_time - h.start_time) ASC,
						h.end_time ASC
				) AS rn

			FROM hasiltryout h
			JOIN users u
				ON u.id_user = h.id_user
			JOIN pesertakelas ps
				ON ps.id_user = h.id_user
			JOIN paketkelas pk
				ON pk.id_paketkelas = ps.id_paketkelas
			WHERE
				h.id_tryout = $1
				AND h.status = 1
				AND h.status_pengerjaan = 'submitted'
		)

		SELECT
			id_user, nama, nama_kelas, nilai, attempt_ke, duration
		FROM ranked_attempt
		WHERE rn = 1
		ORDER BY
			nilai DESC,
			benar DESC,
			duration ASC
	`

	rows, err := r.DB.Query(ctx, query, tryoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []LeaderboardItem

	rank := 1

	for rows.Next() {

		var item LeaderboardItem

		err := rows.Scan(
			&item.UserID,
			&item.Name,
			&item.ClassName,
			&item.Score,
			&item.Attempt,
			&item.Duration,
		)

		if err != nil {
			return nil, err
		}

		item.Rank = rank
		rank++

		result = append(result, item)
	}

	return result, nil
}


//ANCHOR - GET GLOBAL LEADERBOARD SUMMARY
func (r *Repository) GetGlobalLeaderboardSummary(ctx context.Context, tryoutID int) (*LeaderboardSummary, error) {

	query := `
		WITH best_attempt AS (

			SELECT *

			FROM (

				SELECT
					h.id_user,
					h.nilai,

					ROW_NUMBER() OVER (
						PARTITION BY h.id_user
						ORDER BY
							h.nilai DESC,
							h.benar DESC,
							(h.end_time - h.start_time) ASC
					) rn

				FROM hasiltryout h

				WHERE
					h.id_tryout = $1
					AND h.status = 1
					AND h.status_pengerjaan = 'submitted'

			) x

			WHERE rn = 1
		)

		SELECT

			COUNT(*)::INT,
			COALESCE(AVG(nilai),0),
			COALESCE(MAX(nilai),0)

		FROM best_attempt
	`

	var summary LeaderboardSummary

	err := r.DB.QueryRow(ctx, query, tryoutID).Scan(
		&summary.TotalParticipants,
		&summary.AverageScore,
		&summary.HighestScore,
	)

	if err != nil {
		return nil, err
	}

	return &summary, nil
}


//ANCHOR - GET TOTAL ATTEMPT
func (r *Repository) GetTotalAttempt(ctx context.Context, tryoutID int) (int, error) {

	query := `
		SELECT COUNT(*)
		FROM hasiltryout
		WHERE
			id_tryout = $1
			AND status = 1
	`

	var total int

	err := r.DB.QueryRow(ctx, query, tryoutID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}


//ANCHOR - GET USER CLASS ID
func (r *Repository) GetUserClassID(ctx context.Context, userID int) (int, error) {

	query := `
		SELECT id_paketkelas
		FROM pesertakelas
		WHERE
			id_user = $1
			AND status = 1
		LIMIT 1
	`

	var classID int

	err := r.DB.QueryRow(ctx, query, userID).Scan(&classID)
	if err != nil {
		return 0, err
	}

	return classID, nil
}


//ANCHOR - GET CLASS LEADERBOARD
func (r *Repository) GetClassLeaderboard(
	ctx context.Context,
	tryoutID int,
	classID int,
) ([]LeaderboardItem, error) {

	query := `
		WITH ranked_attempt AS (

			SELECT
				h.id_user,
				u.nama,
				pk.nama_kelas,

				h.attempt_ke,
				h.nilai,
				h.benar,
				h.start_time,
				h.end_time,

				EXTRACT(EPOCH FROM (h.end_time - h.start_time))::BIGINT AS duration,

				ROW_NUMBER() OVER (
					PARTITION BY h.id_user
					ORDER BY
						h.nilai DESC,
						h.benar DESC,
						(h.end_time - h.start_time) ASC,
						h.end_time ASC
				) AS rn

			FROM hasiltryout h

			JOIN users u
				ON u.id_user = h.id_user

			JOIN pesertakelas ps
				ON ps.id_user = h.id_user

			JOIN paketkelas pk
				ON pk.id_paketkelas = ps.id_paketkelas

			WHERE
				h.id_tryout = $1
				AND ps.id_paketkelas = $2
				AND h.status = 1
				AND h.status_pengerjaan = 'submitted'
		)

		SELECT
			id_user,
			nama,
			nama_kelas,
			nilai,
			attempt_ke,
			duration

		FROM ranked_attempt

		WHERE rn = 1

		ORDER BY
			nilai DESC,
			benar DESC,
			duration ASC
	`

	rows, err := r.DB.Query(ctx, query, tryoutID, classID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []LeaderboardItem

	rank := 1

	for rows.Next() {

		var item LeaderboardItem

		err := rows.Scan(
			&item.UserID,
			&item.Name,
			&item.ClassName,
			&item.Score,
			&item.Attempt,
			&item.Duration,
		)

		if err != nil {
			return nil, err
		}

		item.Rank = rank
		rank++

		result = append(result, item)
	}

	return result, nil
}


//ANCHOR - GET CLASS LEADERBOARD SUMMARY
func (r *Repository) GetClassLeaderboardSummary(
	ctx context.Context,
	tryoutID int,
	classID int,
) (*LeaderboardSummary, error) {

	query := `
		WITH best_attempt AS (

			SELECT *

			FROM (

				SELECT
					h.id_user,
					h.nilai,

					ROW_NUMBER() OVER (
						PARTITION BY h.id_user
						ORDER BY
							h.nilai DESC,
							h.benar DESC,
							(h.end_time - h.start_time) ASC
					) rn

				FROM hasiltryout h

				JOIN pesertakelas ps
					ON ps.id_user = h.id_user

				WHERE
					h.id_tryout = $1
					AND ps.id_paketkelas = $2
					AND h.status = 1
					AND h.status_pengerjaan = 'submitted'

			) x

			WHERE rn = 1
		)

		SELECT

			COUNT(*)::INT,
			COALESCE(AVG(nilai),0),
			COALESCE(MAX(nilai),0)

		FROM best_attempt
	`

	var summary LeaderboardSummary

	err := r.DB.QueryRow(ctx, query, tryoutID, classID).Scan(
		&summary.TotalParticipants,
		&summary.AverageScore,
		&summary.HighestScore,
	)

	if err != nil {
		return nil, err
	}

	return &summary, nil
}


//ANCHOR - GET CLASS TOTAL ATTEMPT
func (r *Repository) GetClassTotalAttempt(
	ctx context.Context,
	tryoutID int,
	classID int,
) (int, error) {

	query := `
		SELECT COUNT(*)

		FROM hasiltryout h

		JOIN pesertakelas ps
			ON ps.id_user = h.id_user

		WHERE
			h.id_tryout = $1
			AND ps.id_paketkelas = $2
			AND h.status = 1
	`

	var total int

	err := r.DB.QueryRow(ctx, query, tryoutID, classID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
/* ======================== //!SECTION - LEADERBOARD ======================== */



/* ========================================================================== */
/*                 //SECTION - ENDPOINT LIST ANALYTICS SECTION                */
/* ========================================================================== */

// =================================================
// GET STATISTICS OVERVIEW
// =================================================
func (r *Repository) GetStatisticsOverview(
	ctx context.Context,
	tryoutID int,
	classID *int,
) (*StatisticsOverview, error) {

	query := `
		SELECT

			COUNT(*)::INT AS total_attempt,

			COUNT(DISTINCT h.id_user)::INT AS total_participants,

			COUNT(*) FILTER (
				WHERE h.status_pengerjaan = 'submitted'
			)::INT AS completed,

			COUNT(*) FILTER (
				WHERE h.status_pengerjaan <> 'submitted'
			)::INT AS unfinished,

			COALESCE(
				ROUND(
					COUNT(*) FILTER (
						WHERE h.status_pengerjaan = 'submitted'
					)::numeric
					/
					NULLIF(COUNT(*),0)
					*100,
					2
				),
				0
			) AS completion_rate,

			COALESCE(
				ROUND(
					AVG(
						EXTRACT(EPOCH FROM (h.end_time-h.start_time))
					)
				),
				0
			)::BIGINT AS average_duration

		FROM hasiltryout h

		JOIN pesertakelas ps
			ON ps.id_user = h.id_user

		WHERE
			h.id_tryout = $1
			AND h.status = 1
	`

	args := []interface{}{tryoutID}

	if classID != nil {

		query += " AND ps.id_paketkelas = $2"

		args = append(args, *classID)
	}

	var result StatisticsOverview

	err := r.DB.QueryRow(ctx, query, args...).Scan(
		&result.TotalAttempt,
		&result.TotalParticipants,
		&result.Completed,
		&result.Unfinished,
		&result.CompletionRate,
		&result.AverageDuration,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// =================================================
// GET STATISTICS SCORE
// =================================================
func (r *Repository) GetStatisticsScore(
	ctx context.Context,
	tryoutID int,
	classID *int,
) (*StatisticsScore, error) {

	query := `
		WITH best_attempt AS (
			SELECT *
			FROM (
				SELECT
					h.id_user,
					h.nilai,
					h.benar,
					h.start_time,
					h.end_time,
					ROW_NUMBER() OVER (
						PARTITION BY h.id_user
						ORDER BY
							h.nilai DESC,
							h.benar DESC,
							(h.end_time-h.start_time) ASC,
							h.end_time ASC
					) rn
				FROM hasiltryout h
				JOIN pesertakelas ps
					ON ps.id_user = h.id_user
				WHERE
					h.id_tryout = $1
					AND h.status = 1
					AND h.status_pengerjaan='submitted'
	`

	args := []interface{}{tryoutID}

	if classID != nil {

		query += `
					AND ps.id_paketkelas = $2
		`

		args = append(args, *classID)
	}

	query += `
			) x
			WHERE rn=1
		)

		SELECT
			COALESCE(MAX(nilai),0),
			COALESCE(MIN(nilai),0),
			COALESCE(ROUND(AVG(nilai),2),0),
			COALESCE(
				PERCENTILE_CONT(0.5)
				WITHIN GROUP(
					ORDER BY nilai
				),
				0
			)
		FROM best_attempt
	`

	var result StatisticsScore

	err := r.DB.QueryRow(ctx, query, args...).Scan(
		&result.Highest,
		&result.Lowest,
		&result.Average,
		&result.Median,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}


//ANCHOR - GET DISTRIBUTION
func (r *Repository) GetDistribution(
	ctx context.Context,
	tryoutID int,
	classID *int,
) ([]DistributionItem, error) {

	query := `
		WITH best_attempt AS (

			SELECT *

			FROM (

				SELECT

					h.id_user,
					h.nilai,
					h.benar,
					h.start_time,
					h.end_time,

					ROW_NUMBER() OVER (

						PARTITION BY h.id_user

						ORDER BY

							h.nilai DESC,
							h.benar DESC,
							(h.end_time-h.start_time) ASC,
							h.end_time ASC

					) rn

				FROM hasiltryout h

				JOIN pesertakelas ps
					ON ps.id_user = h.id_user

				WHERE

					h.id_tryout = $1
					AND h.status = 1
					AND h.status_pengerjaan='submitted'
	`

	args := []interface{}{tryoutID}

	if classID != nil {

		query += `
			AND ps.id_paketkelas = $2
		`

		args = append(args, *classID)
	}

	query += `
			) x
			WHERE rn = 1
		)
		SELECT nilai
		FROM best_attempt
	`

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// INIT DISTRIBUTION
	distribution := score.EmptyDistribution()

	// COUNT DISTRIBUTION
	for rows.Next() {

		var nilai float64

		if err := rows.Scan(&nilai); err != nil {
			return nil, err
		}

		label := score.GetLabel(nilai)

		distribution[label]++
	}

	// BUILD RESPONSE
	result := make([]DistributionItem, 0, len(score.Labels()))

	for _, label := range score.Labels() {

		result = append(result, DistributionItem{
			Range: label,
			Count: distribution[label],
		})

	}

	return result, nil
}


//ANCHOR - GET QUESTION LIST
func (r *Repository) GetQuestionList(
	ctx context.Context,
	tryoutID int,
) ([]QuestionEntity, error) {

	query := `
		SELECT

			id_soaltryout,
			nomor_urut,
			UPPER(jawaban_benar)

		FROM soaltryout

		WHERE

			id_tryout = $1
			AND status = 1

		ORDER BY nomor_urut
	`

	rows, err := r.DB.Query(ctx, query, tryoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []QuestionEntity

	for rows.Next() {

		var item QuestionEntity

		err := rows.Scan(
			&item.ID,
			&item.Number,
			&item.CorrectAnswer,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}

//? HELPERS
func parseBool(v interface{}) bool {

	switch val := v.(type) {

	case bool:
		return val

	case float64:
		return val == 1

	case int:
		return val == 1

	default:
		return false
	}
} 


//ANCHOR - GET QUESTION ANSWERS
func (r *Repository) GetQuestionAnswers(
	ctx context.Context,
	tryoutID int,
	classID *int,
) ([]QuestionAnswerEntity, int, error) {

	//--------------------------------------------------
	// Mapping nomor soal -> id_soaltryout
	//--------------------------------------------------

	questionNumberMap := make(map[int]int)

	mapQuery := `
		SELECT
			id_soaltryout,
			nomor_urut
		FROM soaltryout
		WHERE
			id_tryout=$1
			AND status=1
	`

	mapRows, err := r.DB.Query(ctx, mapQuery, tryoutID)
	if err != nil {
		return nil, 0, err
	}
	defer mapRows.Close()

	for mapRows.Next() {

		var id int
		var number int

		if err := mapRows.Scan(&id, &number); err != nil {
			return nil, 0, err
		}

		questionNumberMap[number] = id
	}

	//--------------------------------------------------
	// Load jawaban peserta
	//--------------------------------------------------

	query := `
		SELECT
			h.jawaban_user
		FROM hasiltryout h
	`

	args := []interface{}{tryoutID}

	if classID != nil {

		query += `
			JOIN pesertakelas pk
				ON pk.id_user = h.id_user
		`
	}

	query += `
		WHERE

			h.id_tryout = $1
			AND h.status = 1
			AND h.status_pengerjaan='submitted'
	`

	if classID != nil {

		query += `
			AND pk.id_paketkelas=$2
		`

		args = append(args, *classID)
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []QuestionAnswerEntity
	totalParticipant := 0

	for rows.Next() {

		totalParticipant++

		var raw []byte

		if err := rows.Scan(&raw); err != nil {
			return nil, 0, err
		}

		var answers map[string]map[string]interface{}

		if err := json.Unmarshal(raw, &answers); err != nil {
			return nil, 0, err
		}

		for key, value := range answers {

			var questionID int

			//----------------------------------------
			// FORMAT BARU
			//----------------------------------------

			if id, err := strconv.Atoi(key); err == nil {

				questionID = id

			} else {

				//----------------------------------------
				// FORMAT LAMA
				//----------------------------------------

				numberStr := strings.TrimPrefix(key, "soal_")

				number, err := strconv.Atoi(numberStr)
				if err != nil {
					continue
				}

				id, ok := questionNumberMap[number]
				if !ok {
					continue
				}

				questionID = id
			}

			//----------------------------------------
			// answer
			//----------------------------------------

			answer := ""

			if v, ok := value["answer"]; ok {

				if s, ok := v.(string); ok {
					answer = strings.ToUpper(s)
				}
			}

			if answer == "" {

				if v, ok := value["jawaban"]; ok {

					if s, ok := v.(string); ok {
						answer = strings.ToUpper(s)
					}
				}
			}

			//----------------------------------------
			// doubt
			//----------------------------------------

			isDoubt := parseBool(value["ragu"])

			result = append(result, QuestionAnswerEntity{
				QuestionID: questionID,
				Answer:     answer,
				IsDoubt:    isDoubt,
			})
		}
	}

	return result, totalParticipant, nil
}

//ANCHOR - GET QUESTION DETAIL
func (r *Repository) GetQuestionDetail(
	ctx context.Context,
	questionID int,
) (*QuestionDetailEntity, error) {

	query := `
		SELECT
			id_soaltryout,
			nomor_urut,
			UPPER(jawaban_benar)
		FROM soaltryout
		WHERE
			id_soaltryout = $1
			AND status = 1
	`

	var result QuestionDetailEntity

	err := r.DB.QueryRow(ctx, query, questionID).Scan(
		&result.ID,
		&result.Number,
		&result.CorrectAnswer,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


//ANCHOR - GET QUESTION CHOICES
func (r *Repository) GetQuestionChoices(
	ctx context.Context,
	questionID int,
	classID *int,
) ([]QuestionAnswerEntity, int, error) {

	//--------------------------------------------------
	// Cari tryout & mapping nomor -> id
	//--------------------------------------------------

	var tryoutID int

	err := r.DB.QueryRow(ctx,
		`SELECT id_tryout FROM soaltryout WHERE id_soaltryout=$1`,
		questionID,
	).Scan(&tryoutID)

	if err != nil {
		return nil, 0, err
	}

	questionNumberMap := make(map[int]int)

	rowsMap, err := r.DB.Query(ctx, `
		SELECT
			id_soaltryout,
			nomor_urut
		FROM soaltryout
		WHERE
			id_tryout=$1
			AND status=1
	`, tryoutID)
	if err != nil {
		return nil, 0, err
	}
	defer rowsMap.Close()

	for rowsMap.Next() {

		var id, number int

		if err := rowsMap.Scan(&id, &number); err != nil {
			return nil, 0, err
		}

		questionNumberMap[number] = id
	}

	//--------------------------------------------------
	// Load seluruh jawaban
	//--------------------------------------------------

	query := `
		SELECT
			h.jawaban_user
		FROM hasiltryout h
	`

	args := []interface{}{tryoutID}

	if classID != nil {

		query += `
		JOIN pesertakelas pk
			ON pk.id_user = h.id_user
		`
	}

	query += `
		WHERE
			h.id_tryout=$1
			AND h.status=1
			AND h.status_pengerjaan='submitted'
	`

	if classID != nil {

		query += `
			AND pk.id_paketkelas=$2
		`

		args = append(args, *classID)
	}

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []QuestionAnswerEntity
	totalParticipant := 0

	for rows.Next() {

		totalParticipant++

		var raw []byte

		if err := rows.Scan(&raw); err != nil {
			return nil, 0, err
		}

		var answers map[string]map[string]interface{}

		if err := json.Unmarshal(raw, &answers); err != nil {
			return nil, 0, err
		}

		for key, value := range answers {

			var currentQuestionID int

			// ---------- FORMAT BARU ----------
			if id, err := strconv.Atoi(key); err == nil {

				currentQuestionID = id

			} else {

				// ---------- FORMAT LAMA ----------
				numberStr := strings.TrimPrefix(key, "soal_")

				number, err := strconv.Atoi(numberStr)
				if err != nil {
					continue
				}

				id, ok := questionNumberMap[number]
				if !ok {
					continue
				}

				currentQuestionID = id
			}

			// hanya soal yang diminta
			if currentQuestionID != questionID {
				continue
			}

			answer := ""

			if v, ok := value["answer"].(string); ok {
				answer = strings.ToUpper(v)
			}

			if answer == "" {

				if v, ok := value["jawaban"].(string); ok {
					answer = strings.ToUpper(v)
				}
			}

			result = append(result, QuestionAnswerEntity{
				QuestionID: currentQuestionID,
				Answer:     answer,
				IsDoubt:    parseBool(value["ragu"]),
			})
		}
	}

	return result, totalParticipant, nil
}
/* ========================= //!SECTION - ANALYTICS ========================= */