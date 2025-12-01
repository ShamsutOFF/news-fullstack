package dto

type RegisterResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	ID    int64  `json:"id,omitempty"`
}
