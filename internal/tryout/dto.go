package tryout

import "time"

/* ========================================================================== */
/*                  //SECTION - ENDPOINT LIST TRYOUT SECTION                  */
/* ========================================================================== */
type TryoutDTO struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	TotalSoal     int        `json:"total_soal"`
	Duration      int        `json:"duration"`
	MaxAttempt    int        `json:"max_attempt"`
	RemainingAttempt int     `json:"remaining_attempt"`
	AccessStartAt *time.Time `json:"access_start_at,omitempty"`
	AccessEndAt   *time.Time `json:"access_end_at,omitempty"`
	Status        string     `json:"status"` // upcoming | ongoing | ended
}

type StartTryoutResponse struct {
	AttemptToken string `json:"attempt_token"`
	AttemptKe    int    `json:"attempt_ke"`
	Duration     int    `json:"duration"`
	TotalSoal    int    `json:"total_soal"`
	StartTime    string `json:"start_time"`
}

type AttemptEntity struct {
	TryoutID int
	StartTime time.Time
	EndTime   *time.Time
	Duration int
	Status   string
}

type OngoingAttemptEntity struct {
	ID           int
	AttemptToken string
	AttemptKe    int
	StartTime    time.Time
	EndTime      *time.Time
}

type SoalDTO struct {
	ID         int               `json:"id"`
	Nomor      int               `json:"nomor"`
	Pertanyaan string            `json:"pertanyaan"`
	Pilihan    map[string]string `json:"pilihan"`
}

type GetSoalResponse struct {
	AttemptToken string     `json:"attempt_token"`
	Duration     int        `json:"duration"`
	Remaining    int64      `json:"remaining_time"`
	Questions    []SoalDTO  `json:"questions"`
}

type AnswerPayload struct {
	Answer string `json:"answer"`
	Ragu   bool   `json:"ragu"`
}

type SaveAnswerRequest struct {
	Answers map[string]AnswerPayload `json:"answers"`
}

type SubmitResponse struct {
	Score  float64 `json:"score"`
	Benar  int     `json:"benar"`
	Salah  int     `json:"salah"`
	Kosong int     `json:"kosong"`
	RaguRagu int `json:"ragu_ragu"`
}

type ResumeResponse struct {
	AttemptToken string             	  `json:"attempt_token"`
	TryoutID     int                	  `json:"tryout_id"`
	Duration     int            	      `json:"duration"`
	Remaining    int64      	          `json:"remaining_time"`
	Status       string		              `json:"status"`
	Answers		 map[string]AnswerPayload `json:"answers"`
	Questions    []SoalDTO           	  `json:"questions"`
}
/* =========================== //!SECTION - TRYOUT ========================== */



/* ========================================================================== */
/*                  //SECTION - ENDPOINT LIST REPORT SECTION                  */
/* ========================================================================== */
type TryoutReportDTO struct {
	AttemptToken string    `json:"attempt_token"`
	TryoutID     int       `json:"tryout_id"`
	Title        string    `json:"title"`
	AttemptKe    int       `json:"attempt_ke"`
	Score        float64   `json:"score"`
	Benar        int       `json:"benar"`
	Salah        int       `json:"salah"`
	Kosong       int       `json:"kosong"`
	RaguRagu	 int		`json:"ragu_ragu"`
	Tanggal      time.Time `json:"tanggal"`
}

type ReviewDTO struct {
	ID            int               `json:"id"`
	Nomor         int               `json:"nomor"`
	Pertanyaan    string            `json:"pertanyaan"`
	Pilihan       map[string]string `json:"pilihan"`
	UserAnswer    string            `json:"user_answer"`
	CorrectAnswer string            `json:"correct_answer"`
	Pembahasan    string            `json:"pembahasan"`
	Status        string            `json:"status"`
	IsRagu        bool              `json:"is_ragu"`
}
/* =========================== //!SECTION - REPORT ========================== */



/* ========================================================================== */
/*                //SECTION - ENDPOINT LIST LEADERBOARD SECTION               */
/* ========================================================================== */
type LeaderboardItem struct {
	Rank      int     `json:"rank"`
	UserID    int     `json:"user_id"`
	Name      string  `json:"name"`
	ClassName string  `json:"class"`
	Score     float64 `json:"score"`
	Attempt   int     `json:"attempt"`
	Duration  int64   `json:"duration"`
}


type LeaderboardSummary struct {
    TotalParticipants int     `json:"total_participants"`
    TotalAttempt      int     `json:"total_attempt"`
    AverageScore      float64 `json:"average_score"`
    HighestScore      float64 `json:"highest_score"`
    MyScore           float64 `json:"my_score"`
    MyRank            int     `json:"my_rank"`
}

type LeaderboardResponse struct {
    Summary    LeaderboardSummary `json:"summary"`
    Leaderboard []LeaderboardItem `json:"leaderboard"`
}
/* ======================== //!SECTION - LEADERBOARD ======================== */



/* ========================================================================== */
/*                //SECTION - ENDPOINT LIST ANALYTICS SECTION                */
/* ========================================================================== */
type StatisticsOverview struct {
	TotalParticipants int     `json:"total_participants"`
	TotalAttempt      int     `json:"total_attempt"`

	Completed         int     `json:"completed"`
	Unfinished        int     `json:"unfinished"`

	CompletionRate    float64 `json:"completion_rate"`
	AverageDuration   int64   `json:"average_duration"`
}

type StatisticsScore struct {
	Highest float64 `json:"highest"`
	Lowest  float64 `json:"lowest"`
	Average float64 `json:"average"`
	Median  float64 `json:"median"`
}

type StatisticsCompletion struct {
	CompletionRate      float64 `json:"completion_rate"`
	AverageDuration     int64   `json:"average_duration"`
}

type StatisticsResponse struct {
	Overview StatisticsOverview `json:"overview"`
	Score    StatisticsScore    `json:"score"`
}

type DistributionItem struct {
	Range string `json:"range"`
	Count int    `json:"count"`
}

type DistributionResponse struct {
	Distribution []DistributionItem `json:"distribution"`
}

type QuestionAnalysisDTO struct {
	ID     int `json:"id"`
	Number int `json:"number"`

	Correct int `json:"correct"`
	Wrong   int `json:"wrong"`
	Blank   int `json:"blank"`
	Doubt   int `json:"doubt"`

	CorrectRate float64 `json:"correct_rate"`
}

type QuestionAnalysisEntity struct {
	ID            int
	Number        int
	CorrectAnswer string

	Correct int
	Wrong   int
	Blank   int
	Doubt   int

	TotalParticipant int
}

type QuestionEntity struct {
	ID            int
	Number        int
	CorrectAnswer string
}

type QuestionAnswerEntity struct {
	QuestionID int
	Answer     string
	IsDoubt    bool
}

type QuestionChoiceDTO struct {
	ID     int    `json:"id"`
	Number int    `json:"number"`

	CorrectAnswer string `json:"correct_answer"`

	A int `json:"a"`
	B int `json:"b"`
	C int `json:"c"`
	D int `json:"d"`
	E int `json:"e"`

	Blank int `json:"blank"`
	Doubt int `json:"doubt"`

	TotalParticipant int `json:"total_participant"`
	TotalAnswer      int `json:"total_answer"`
}

type QuestionDetailEntity struct {
	ID            int
	Number        int
	CorrectAnswer string
}
/* ======================== //!SECTION - ANALYTICS ======================== */