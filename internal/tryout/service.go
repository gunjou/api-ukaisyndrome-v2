package tryout

import (
	"api-ukaisyndrome-v2/pkg/timeutil"
	"context"
	"encoding/json"
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	Repo *Repository
}


/* ========================================================================== */
/*                  //SECTION - ENDPOINT LIST TRYOUT SECTION                  */
/* ========================================================================== */

//ANCHOR - GET TRYOUT BY USER
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


//ANCHOR - RESOLVE EXPIRED ATTEMPTS
func (s *Service) ResolveExpiredAttempts(
	ctx context.Context,
	userID int,
	tryoutID int,
) error {

	attempts, err := s.Repo.GetOngoingAttempts(
		ctx,
		userID,
		tryoutID,
	)
	if err != nil {
		return err
	}

	now := timeutil.Now()

	for _, attempt := range attempts {

		// jika end_time belum ada (data lama), skip dulu
		if attempt.EndTime == nil {
			continue
		}

		// masih aktif
		if now.Before(*attempt.EndTime) {
			continue
		}

		nilai, benar, salah, kosong, ragu, err :=
			s.CalculateAttemptResult(
				ctx,
				tryoutID,
				attempt.AttemptToken,
				userID,
			)

		if err != nil {
			return err
		}

		err = s.Repo.UpdateSubmittedAttempt(
			ctx,
			attempt.ID,
			nilai,
			benar,
			salah,
			kosong,
			ragu,
		)
		if err != nil {
			return err
		}
	}

	return nil
}


//ANCHOR - START TRYOUT
func (s *Service) StartTryout(
	ctx context.Context,
	userID int,
	tryoutID int,
) (*StartTryoutResponse, error) {

	//--------------------------------------------------
	// resolve seluruh attempt yang sudah expired
	//--------------------------------------------------

	err := s.ResolveExpiredAttempts(
		ctx,
		userID,
		tryoutID,
	)
	if err != nil {
		return nil, err
	}

	//--------------------------------------------------
	// validasi akses
	//--------------------------------------------------

	t, err := s.Repo.GetTryoutForStart(
		ctx,
		userID,
		tryoutID,
	)
	if err != nil {
		return nil, errors.New("tryout not accessible")
	}

	now := timeutil.Now()

	//--------------------------------------------------
	// validasi waktu
	//--------------------------------------------------

	if t.AccessStartAt != nil && now.Before(*t.AccessStartAt) {
		return nil, errors.New("tryout not started")
	}

	if t.AccessEndAt != nil && now.After(*t.AccessEndAt) {
		return nil, errors.New("tryout expired")
	}

	//--------------------------------------------------
	// validasi attempt
	//--------------------------------------------------

	count, err := s.Repo.CountAttempt(
		ctx,
		userID,
		tryoutID,
	)
	if err != nil {
		return nil, err
	}

	if t.MaxAttempt > 0 && count >= t.MaxAttempt {
		return nil, errors.New("max attempt reached")
	}

	//--------------------------------------------------
	// generate attempt
	//--------------------------------------------------

	attemptToken := uuid.New().String()
	attemptKe := count + 1

	err = s.Repo.InsertAttempt(
		ctx,
		userID,
		tryoutID,
		attemptToken,
		attemptKe,
		t.Duration,
	)
	if err != nil {
		return nil, err
	}

	return &StartTryoutResponse{
		AttemptToken: attemptToken,
		AttemptKe:    attemptKe,
		Duration:     t.Duration,
		TotalSoal:    t.TotalSoal,
		StartTime:    now.Format(time.RFC3339),
	}, nil
}


//ANCHOR - GET SOAL TRYOUT
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


//ANCHOR - SAVE ANSWER (FOR AUTOSAVE)
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
	// if len(answers) == 0 {
	// 	return errors.New("answers cannot be empty")
	// }

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


//ANCHOR - CALCULATE ATTEMPT RESULT
func (s *Service) CalculateAttemptResult(ctx context.Context, tryoutID int, attemptToken string, userID int,
) (
	float64,
	int,
	int,
	int,
	int,
	error,
) {

	// answer key
	keys, err := s.Repo.GetAnswerKey(ctx, tryoutID)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	// user answer
	userAnswers, err := s.Repo.GetUserAnswers(
		ctx,
		attemptToken,
		userID,
	)
	if err != nil {
		return 0, 0, 0, 0, 0, err
	}

	var benar, salah, kosong, ragu int

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
			ragu++
		}
	}

	total := len(keys)
	nilai := 0.0
	if total > 0 {
		nilai = float64(benar) / float64(total) * 100
	}

	return nilai, benar, salah, kosong, ragu, nil
}


//ANCHOR - SUBMIT TRYOUT
func (s *Service) SubmitTryout(ctx context.Context, userID int, attemptToken string) (*SubmitResponse, error) {

	// 🔐 VALIDASI ATTEMPT
	attempt, err := s.Repo.GetAttempt(ctx, attemptToken, userID)
	if err != nil {
		return nil, errors.New("invalid attempt")
	}

	if attempt.Status == "submitted" {
		return nil, errors.New("already submitted")
	}

	nilai, benar, salah, kosong, raguRagu, err :=
		s.CalculateAttemptResult(
			ctx,
			attempt.TryoutID,
			attemptToken,
			userID,
		)

	if err != nil {
		return nil, err
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



//ANCHOR - RESUME ATTEMPT
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

//ANCHOR - GET ONGOING TRYOUT
func (s *Service) GetOngoingTryout(
	ctx context.Context,
	userID int,
) (*OngoingTryoutResponse, error) {

	attempts, err := s.Repo.GetOngoingTryout(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := timeutil.Now()

	response := &OngoingTryoutResponse{
		Ongoing: make([]OngoingTryoutDTO, 0),
		Expired: make([]OngoingTryoutDTO, 0),
	}

	for _, item := range attempts {

		var answers map[string]interface{}

		if len(item.JawabanUser) > 0 {

			if err := json.Unmarshal(item.JawabanUser, &answers); err != nil {
				answers = map[string]interface{}{}
			}

		} else {
			answers = map[string]interface{}{}
		}

		dto := OngoingTryoutDTO{
			IDHasilTryout:    item.IDHasilTryout,
			IDTryout:         item.IDTryout,
			AttemptToken:     item.AttemptToken,
			StartTime:        item.StartTime,
			EndTime:          item.EndTime,
			JawabanUser:      answers,
			StatusPengerjaan: item.StatusPengerjaan,
		}

		// TANPA END TIME
		if item.EndTime == nil {

			response.Ongoing = append(response.Ongoing, dto)
			continue
		}

		// EXPIRED
		if now.After(*item.EndTime) {

			response.Expired = append(response.Expired, dto)
			continue
		}

		// MASIH BERJALAN
		response.Ongoing = append(response.Ongoing, dto)
	}

	return response, nil
}
/* =========================== //!SECTION - TRYOUT ========================== */



/* ========================================================================== */
/*                  //SECTION - ENDPOINT LIST REPORT SECTION                  */
/* ========================================================================== */

//ANCHOR - GET TRYOUT REPORTS
func (s *Service) GetReports(ctx context.Context, userID int) ([]TryoutReportDTO, error) {
	return s.Repo.GetUserReports(ctx, userID)
}

//ANCHOR - GET TRYOUT REVIEW
func (s *Service) GetReview(ctx context.Context, userID int, token string) ([]ReviewDTO, error) {
	return s.Repo.GetReview(ctx, token, userID)
}
/* =========================== //!SECTION - REPORT ========================== */




/* ========================================================================== */
/*                //SECTION - ENDPOINT LIST LEADERBOARD SECTION               */
/* ========================================================================== */

//ANCHOR - GET GLOBAL LEADERBOARD
func (s *Service) GetGlobalLeaderboard(
	ctx context.Context,
	userID int,
	tryoutID int,
) (*LeaderboardResponse, error) {

	// Summary
	summary, err := s.Repo.GetGlobalLeaderboardSummary(ctx, tryoutID)
	if err != nil {
		return nil, err
	}

	// Total Attempt
	totalAttempt, err := s.Repo.GetTotalAttempt(ctx, tryoutID)
	if err != nil {
		return nil, err
	}

	summary.TotalAttempt = totalAttempt

	// Leaderboard
	leaderboard, err := s.Repo.GetGlobalLeaderboard(ctx, tryoutID)
	if err != nil {
		return nil, err
	}

	// Cari ranking & score user login
	for _, item := range leaderboard {

		if item.UserID == userID {

			summary.MyRank = item.Rank
			summary.MyScore = item.Score

			break
		}
	}

	return &LeaderboardResponse{
		Summary:     *summary,
		Leaderboard: leaderboard,
	}, nil
}


//ANCHOR - GET CLASS LEADERBOARD
func (s *Service) GetClassLeaderboard(
	ctx context.Context,
	userID int,
	role string,
	tryoutID int,
	classID int,
) (*LeaderboardResponse, error) {

	// Tentukan kelas yang digunakan
	switch role {

	case "admin", "mentor":

		if classID == 0 {
			return nil, errors.New("class_id is required")
		}

	case "peserta":

		var err error

		classID, err = s.Repo.GetUserClassID(ctx, userID)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("invalid role")
	}

	// Summary
	summary, err := s.Repo.GetClassLeaderboardSummary(
		ctx,
		tryoutID,
		classID,
	)
	if err != nil {
		return nil, err
	}

	// Total Attempt
	totalAttempt, err := s.Repo.GetClassTotalAttempt(
		ctx,
		tryoutID,
		classID,
	)
	if err != nil {
		return nil, err
	}

	summary.TotalAttempt = totalAttempt

	// Leaderboard
	leaderboard, err := s.Repo.GetClassLeaderboard(
		ctx,
		tryoutID,
		classID,
	)
	if err != nil {
		return nil, err
	}

	// My Rank
	for _, item := range leaderboard {

		if item.UserID == userID {

			summary.MyRank = item.Rank
			summary.MyScore = item.Score

			break
		}
	}

	return &LeaderboardResponse{
		Summary:     *summary,
		Leaderboard: leaderboard,
	}, nil
}
/* ======================== //!SECTION - LEADERBOARD ======================== */



/* ========================================================================== */
/*                 //SECTION - ENDPOINT LIST ANALYTICS SECTION                */
/* ========================================================================== */

// GET STATISTICS
func (s *Service) GetStatistics(
	ctx context.Context,
	userID int,
	role string,
	tryoutID int,
	classID int,
) (*StatisticsResponse, error) {

	// Tentukan scope
	var classFilter *int

	switch role {

	case "admin", "mentor":

		// class_id optional
		if classID != 0 {
			classFilter = &classID
		}

	case "peserta":

		userClassID, err := s.Repo.GetUserClassID(ctx, userID)
		if err != nil {
			return nil, err
		}

		classFilter = &userClassID

	default:
		return nil, errors.New("invalid role")
	}

	// Overview
	overview, err := s.Repo.GetStatisticsOverview(
		ctx,
		tryoutID,
		classFilter,
	)
	if err != nil {
		return nil, err
	}

	// Score
	score, err := s.Repo.GetStatisticsScore(
		ctx,
		tryoutID,
		classFilter,
	)
	if err != nil {
		return nil, err
	}

	return &StatisticsResponse{
		Overview: *overview,
		Score:    *score,
	}, nil
}

//ANCHOR - GET DISTRIBUTION
func (s *Service) GetDistribution(
	ctx context.Context,
	userID int,
	role string,
	tryoutID int,
	classID int,
) (*DistributionResponse, error) {

	var classFilter *int

	switch role {

	case "admin", "mentor":

		if classID != 0 {
			classFilter = &classID
		}

	case "peserta":

		userClassID, err := s.Repo.GetUserClassID(ctx, userID)
		if err != nil {
			return nil, err
		}

		classFilter = &userClassID

	default:
		return nil, errors.New("invalid role")
	}

	data, err := s.Repo.GetDistribution(
		ctx,
		tryoutID,
		classFilter,
	)
	if err != nil {
		return nil, err
	}

	return &DistributionResponse{
		Distribution: data,
	}, nil
}


// =================================================
// GET QUESTION ANALYSIS
// =================================================
func (s *Service) GetQuestionAnalysis(
	ctx context.Context,
	userID int,
	role string,
	tryoutID int,
	classID int,
	sortBy string,
) ([]QuestionAnalysisDTO, error) {

	//--------------------------------------------------
	// DETERMINE CLASS FILTER
	//--------------------------------------------------

	var classFilter *int

	switch role {

	case "admin", "mentor":

		if classID != 0 {
			classFilter = &classID
		}

	case "peserta":

		userClassID, err := s.Repo.GetUserClassID(ctx, userID)
		if err != nil {
			return nil, err
		}

		classFilter = &userClassID

	default:
		return nil, errors.New("invalid role")
	}

	//--------------------------------------------------
	// LOAD DATA
	//--------------------------------------------------

	questions, err := s.Repo.GetQuestionList(ctx, tryoutID)
	if err != nil {
		return nil, err
	}

	answers, totalParticipant, err := s.Repo.GetQuestionAnswers(
		ctx,
		tryoutID,
		classFilter,
	)
	if err != nil {
		return nil, err
	}

	//--------------------------------------------------
	// INIT MAP
	//--------------------------------------------------

	resultMap := make(map[int]*QuestionAnalysisDTO)
	correctAnswerMap := make(map[int]string)

	for _, q := range questions {

		correctAnswerMap[q.ID] = strings.ToUpper(q.CorrectAnswer)

		resultMap[q.ID] = &QuestionAnalysisDTO{
			ID:      q.ID,
			Number:  q.Number,
			Correct: 0,
			Wrong:   0,
			Blank:   totalParticipant,
			Doubt:   0,
		}
	}

	//--------------------------------------------------
	// PROCESS ANSWERS
	//--------------------------------------------------

	for _, ans := range answers {

		stat, ok := resultMap[ans.QuestionID]
		if !ok {
			continue
		}

		// Blank hanya berkurang jika ada jawaban
		if ans.Answer != "" {
			stat.Blank--
		}

		if ans.IsDoubt {
			stat.Doubt++
		}

		if ans.Answer == "" {
			continue
		}

		if ans.Answer == correctAnswerMap[ans.QuestionID] {
			stat.Correct++
		} else {
			stat.Wrong++
		}
	}

	//--------------------------------------------------
	// BUILD RESPONSE
	//--------------------------------------------------

	result := make([]QuestionAnalysisDTO, 0, len(questions))

	for _, q := range questions {

		stat := resultMap[q.ID]

		if totalParticipant > 0 {
			stat.CorrectRate = math.Round(
				(float64(stat.Correct)/float64(totalParticipant))*100,
			) / 100
		}

		result = append(result, *stat)
	}

	//--------------------------------------------------
	// SORT
	//--------------------------------------------------

	switch strings.ToLower(sortBy) {

	case "hardest":

		sort.Slice(result, func(i, j int) bool {
			return result[i].CorrectRate < result[j].CorrectRate
		})

	case "easiest":

		sort.Slice(result, func(i, j int) bool {
			return result[i].CorrectRate > result[j].CorrectRate
		})

	default:

		sort.Slice(result, func(i, j int) bool {
			return result[i].Number < result[j].Number
		})
	}

	return result, nil
}


//ANCHOR - GET QUESTION CHOICE ANALYSIS
func (s *Service) GetQuestionChoices(
	ctx context.Context,
	userID int,
	role string,
	questionID int,
	classID int,
) (*QuestionChoiceDTO, error) {

	//--------------------------------------------------
	// DETERMINE CLASS FILTER
	//--------------------------------------------------

	var classFilter *int

	switch role {

	case "admin", "mentor":

		if classID != 0 {
			classFilter = &classID
		}

	case "peserta":

		userClassID, err := s.Repo.GetUserClassID(ctx, userID)
		if err != nil {
			return nil, err
		}

		classFilter = &userClassID

	default:
		return nil, errors.New("invalid role")
	}

	//--------------------------------------------------
	// LOAD DATA
	//--------------------------------------------------

	question, err := s.Repo.GetQuestionDetail(ctx, questionID)
	if err != nil {
		return nil, err
	}

	answers, totalParticipant, err := s.Repo.GetQuestionChoices(
		ctx,
		questionID,
		classFilter,
	)
	if err != nil {
		return nil, err
	}

	//--------------------------------------------------
	// BUILD RESULT
	//--------------------------------------------------

	result := &QuestionChoiceDTO{
		ID:               question.ID,
		Number:           question.Number,
		CorrectAnswer:    question.CorrectAnswer,
		TotalParticipant: totalParticipant,
	}

	//--------------------------------------------------
	// COUNT ANSWERS
	//--------------------------------------------------

	for _, ans := range answers {

		if ans.IsDoubt {
			result.Doubt++
		}

		switch ans.Answer {

		case "A":
			result.A++
			result.TotalAnswer++

		case "B":
			result.B++
			result.TotalAnswer++

		case "C":
			result.C++
			result.TotalAnswer++

		case "D":
			result.D++
			result.TotalAnswer++

		case "E":
			result.E++
			result.TotalAnswer++

		default:
			result.Blank++
		}
	}

	//--------------------------------------------------
	// HANDLE LEGACY DATA
	//--------------------------------------------------

	// Jika repository tidak mengembalikan blank (karena format lama),
	// pastikan blank minimal sesuai selisih peserta dan jawaban.
	if result.Blank == 0 && result.TotalParticipant > result.TotalAnswer {
		result.Blank = result.TotalParticipant - result.TotalAnswer
	}

	return result, nil
}
/* ========================= //!SECTION - ANALYTICS ========================= */