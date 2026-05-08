package tryout

import "time"

type TryoutDTO struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	TotalSoal     int        `json:"total_soal"`
	Duration      int        `json:"duration"`
	MaxAttempt    int        `json:"max_attempt"`
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
	Duration int
	Status   string
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