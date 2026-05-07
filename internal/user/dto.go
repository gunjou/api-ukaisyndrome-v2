package user

type ClassDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	BatchID     int    `json:"id_batch"`
	BatchName   string `json:"batch"`
	StatusBatch int    `json:"status_batch"`
}

type MentorshipDTO struct {
	ID          int    `json:"id"`
	MentorID   int    `json:"mentor_id"`
	MentorshipName string `json:"mentorship_name"`
	MentorName  string `json:"mentor_name"`
}

type MeResponse struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Role     string     `json:"role"`
	Mentorships []MentorshipDTO `json:"mentorships"`

	// optional
	Nickname *string    `json:"nickname,omitempty"`
	Classes  []ClassDTO `json:"classes,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}