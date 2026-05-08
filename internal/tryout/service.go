package tryout

import (
	"api-ukaisyndrome-v2/pkg/timeutil"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	Repo *Repository
}


// =================================================
// GET TRYOUT BY USER
// =================================================
func (s *Service) GetTryoutPeserta(ctx context.Context, userID int) ([]TryoutDTO, error) {

	data, err := s.Repo.GetTryoutByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := timeutil.Now()

	for i := range data {

		start := data[i].AccessStartAt
		end := data[i].AccessEndAt

		if start != nil && now.Before(*start) {
			data[i].Status = "upcoming"
			continue
		}

		if start != nil && end != nil && now.After(*start) && now.Before(*end) {
			data[i].Status = "ongoing"
			continue
		}

		if end != nil && now.After(*end) {
			data[i].Status = "ended"
			continue
		}

		data[i].Status = "unknown"
	}

	return data, nil
}


// =================================================
// GET TRYOUT FOR START
// =================================================
func (s *Service) StartTryout(ctx context.Context, userID int, tryoutID int) (*StartTryoutResponse, error) {

	// 1. validasi akses
	t, err := s.Repo.GetTryoutForStart(ctx, userID, tryoutID)
	if err != nil {
		return nil, errors.New("tryout not accessible")
	}

	now := timeutil.Now()

	// 2. validasi waktu
	if t.AccessStartAt != nil && now.Before(*t.AccessStartAt) {
		return nil, errors.New("tryout not started")
	}

	if t.AccessEndAt != nil && now.After(*t.AccessEndAt) {
		return nil, errors.New("tryout expired")
	}

	// 3. validasi attempt
	count, err := s.Repo.CountAttempt(ctx, userID, tryoutID)
	if err != nil {
		return nil, err
	}

	if t.MaxAttempt > 0 && count >= t.MaxAttempt {
		return nil, errors.New("max attempt reached")
	}

	// 4. generate token
	attemptToken := uuid.New().String()
	attemptKe := count + 1

	// 5. insert
	err = s.Repo.InsertAttempt(ctx, userID, tryoutID, attemptToken, attemptKe)
	if err != nil {
		return nil, err
	}

	return &StartTryoutResponse{
		AttemptToken: attemptToken,
		AttemptKe:    attemptKe,
		Duration:     t.Duration,
		TotalSoal:    t.TotalSoal,
		StartTime:    timeutil.Now().Format(time.RFC3339),
	}, nil
}


// =================================================
// GET SOAL TRYOUT
// =================================================
func (s *Service) GetSoalTryout(ctx context.Context, userID int, attemptToken string) (*GetSoalResponse, error) {

	// 🔐 VALIDASI ATTEMPT
	attempt, err := s.Repo.GetAttempt(ctx, attemptToken, userID)
	if err != nil {
		return nil, errors.New("invalid attempt")
	}

	if attempt.Status != "ongoing" {
		return nil, errors.New("tryout already finished")
	}

	// ⏱ VALIDASI TIMER
	expiredAt := attempt.StartTime.Add(time.Duration(attempt.Duration) * time.Minute)

	now := timeutil.Now()

	if now.After(expiredAt) {
		return nil, errors.New("time is up")
	}

	// 📦 ambil soal
	questions, err := s.Repo.GetSoalByTryout(ctx, attempt.TryoutID)
	if err != nil {
		return nil, err
	}

	remaining := int64(time.Until(expiredAt).Seconds())

	return &GetSoalResponse{
		AttemptToken: attemptToken,
		Duration:     attempt.Duration,
		Remaining:    remaining,
		Questions:    questions,
	}, nil
}


// =================================================
// SAVE ANSWER (FOR AUTOSAVE)
// =================================================
func (s *Service) SaveAnswers(
	ctx context.Context,
	userID int,
	attemptToken string,
	answers map[string]AnswerPayload,
) error {

	// 🔐 VALIDASI ATTEMPT
	attempt, err := s.Repo.GetAttempt(ctx, attemptToken, userID)
	if err != nil {
		return errors.New("invalid attempt")
	}

	if attempt.Status != "ongoing" {
		return errors.New("tryout already finished")
	}

	// ⏱ VALIDASI WAKTU
	expiredAt := attempt.StartTime.Add(time.Duration(attempt.Duration) * time.Minute)
	if time.Until(expiredAt) <= 0 {
		return errors.New("time is up")
	}

	// 🔥 VALIDASI INPUT
	if len(answers) == 0 {
		return errors.New("answers cannot be empty")
	}

	// validasi pilihan + ragu
	for _, v := range answers {
		if v.Answer != "" {
			switch v.Answer {
			case "A", "B", "C", "D", "E":
				// ok
			default:
				return errors.New("invalid answer value")
			}
		}
	}

	return s.Repo.SaveAnswers(ctx, attemptToken, userID, answers)
}


// =================================================
// SUBMIT TRYOUT
// =================================================
func (s *Service) SubmitTryout(ctx context.Context, userID int, attemptToken string) (*SubmitResponse, error) {

	// 🔐 VALIDASI ATTEMPT
	attempt, err := s.Repo.GetAttempt(ctx, attemptToken, userID)
	if err != nil {
		return nil, errors.New("invalid attempt")
	}

	if attempt.Status == "submitted" {
		return nil, errors.New("already submitted")
	}

	// 🔑 ambil answer key
	keys, err := s.Repo.GetAnswerKey(ctx, attempt.TryoutID)
	if err != nil {
		return nil, err
	}

	// 👤 ambil jawaban user (sudah struct {answer, ragu})
	userAnswers, err := s.Repo.GetUserAnswers(ctx, attemptToken, userID)
	if err != nil {
		return nil, err
	}

	var benar, salah, kosong, raguRagu int

	for id, correct := range keys {

		ua, ok := userAnswers[strconv.Itoa(id)]

		if !ok || ua.Answer == "" {
			kosong++
			continue
		}

		if ua.Answer == correct {
			benar++
		} else {
			salah++
		}

		if ua.Ragu {
			raguRagu++
		}
	}

	total := len(keys)
	nilai := 0.0
	if total > 0 {
		nilai = float64(benar) / float64(total) * 100
	}

	// 🔥 update DB (sudah include ragu)
	if err := s.Repo.SubmitResult(
		ctx,
		attemptToken,
		userID,
		benar,
		salah,
		kosong,
		raguRagu,
		nilai,
	); err != nil {
		return nil, err
	}

	return &SubmitResponse{
		Score:     nilai,
		Benar:     benar,
		Salah:     salah,
		Kosong:    kosong,
		RaguRagu:  raguRagu,
	}, nil
}



// =================================================
// RESUME ATTEMPT
// =================================================
func (s *Service) ResumeAttempt(ctx context.Context, userID int, attemptToken string) (*ResumeResponse, error) {

	// 🔐 VALIDASI ATTEMPT
	attempt, err := s.Repo.GetAttempt(ctx, attemptToken, userID)
	if err != nil {
		return nil, errors.New("invalid attempt")
	}

	if attempt.Status == "submitted" {
		return nil, errors.New("tryout already submitted")
	}

	// ⏱ HITUNG SISA WAKTU
	expiredAt := attempt.StartTime.Add(time.Duration(attempt.Duration) * time.Minute)

	remaining := int64(time.Until(expiredAt).Seconds())
	if remaining < 0 {
		remaining = 0
	}

	// 📦 ambil soal
	questions, err := s.Repo.GetSoalByTryout(ctx, attempt.TryoutID)
	if err != nil {
		return nil, err
	}

	// 📦 ambil jawaban user (struct baru)
	answers, err := s.Repo.GetUserAnswers(ctx, attemptToken, userID)
	if err != nil {
		return nil, err
	}

	return &ResumeResponse{
		AttemptToken: attemptToken,
		TryoutID:     attempt.TryoutID,
		Duration:     attempt.Duration,
		Remaining:    remaining,
		Status:       attempt.Status,
		Answers:      answers, // <- penting: map[string]AnswerPayload
		Questions:    questions,
	}, nil
}


// =================================================
// GET TRYOUT REPORTS
// =================================================
func (s *Service) GetReports(ctx context.Context, userID int) ([]TryoutReportDTO, error) {
	return s.Repo.GetUserReports(ctx, userID)
}


// =================================================
// GET TRYOUT REVIEW
// =================================================
func (s *Service) GetReview(ctx context.Context, userID int, token string) ([]ReviewDTO, error) {
	return s.Repo.GetReview(ctx, token, userID)
}