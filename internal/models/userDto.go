package models

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginToken struct {
	ID        uint
	Email     string
	Role      string
	FirstName string
	LastName  string
}

type LoginUserResponse struct {
	ID            uint           `json:"id"`
	Email         string         `json:"email"`
	Role          string         `json:"role"`
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	PhotoUrl      string         `json:"photo_url"`
	StudentStatus *StudentStatus `json:"student_status,omitempty"`
}

func (l *LoginUserResponse) CreateLoginToken() *LoginToken {
	return &LoginToken{
		ID:        l.ID,
		Email:     l.Email,
		Role:      l.Role,
		FirstName: l.FirstName,
		LastName:  l.LastName,
	}
}
