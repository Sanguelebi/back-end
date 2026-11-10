package models

type LoginToken struct {
	ID        uint
	Email     string
	Role      string
	FirstName string
	LastName  string
}

type RegisterUserRequest struct {
	Email         string         `json:"email"`
	Password      string         `json:"password"`
	Role          Role           `json:"role"`
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	Gender        Gender         `json:"gender"`
	PersonalNo    string         `json:"personal_no"`
	PhoneNumber   string         `json:"phone_number"`
	PhotoUrl      string         `json:"photo_url"`
	StudentStatus *StudentStatus `json:"student_status"`
}

func (r *RegisterUserRequest) CreateTableUser() *Users {
	return &Users{
		Email:         r.Email,
		Password:      r.Password,
		Role:          r.Role,
		FirstName:     r.FirstName,
		LastName:      r.LastName,
		Gender:        r.Gender,
		PersonalNo:    r.PersonalNo,
		PhoneNumber:   r.PhoneNumber,
		PhotoUrl:      r.PhotoUrl,
		StudentStatus: r.StudentStatus,
	}
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
