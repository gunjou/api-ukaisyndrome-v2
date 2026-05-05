package user

type ClassDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	BatchID     int    `json:"id_batch"`
	BatchName   string `json:"batch"`
	StatusBatch int    `json:"status_batch"`
}

type MeResponse struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Role     string     `json:"role"`

	// optional
	Nickname *string    `json:"nickname,omitempty"`
	Classes  []ClassDTO `json:"classes,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}