package models

type LoginToken struct {
	ID        uint
	Email     string
	Role      string
	FirstName string
	LastName  string
}

// INFO: Required attribute adds additional security.
type RegisterUserRequest struct {
	Email         string         `json:"email" validate:"required,email"`
	Password      string         `json:"password" validate:"required,min=8"`
	Role          Role           `json:"role"`
	FirstName     string         `json:"first_name" validate:"required"`
	LastName      string         `json:"last_name" validate:"required"`
	Gender        Gender         `json:"gender"`
	PersonalNo    string         `json:"personal_no"`
	PhoneNumber   string         `json:"phone_number"`
	PhotoUrl      string         `json:"photo_url"`
	StudentStatus *StudentStatus `json:"student_status,omitempty"`
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

// INFO: Only dirtyFields will be sent over from Front-End.
type EditUserRequest struct {
	Email         *string        `json:"email,omitempty"`
	Role          *Role          `json:"role,omitempty"`
	FirstName     *string        `json:"first_name,omitempty"`
	LastName      *string        `json:"last_name,omitempty"`
	Gender        *Gender        `json:"gender,omitempty"`
	PersonalNo    *string        `json:"personal_no,omitempty"`
	PhoneNumber   *string        `json:"phone_number,omitempty"`
	PhotoUrl      *string        `json:"photo_url,omitempty"`
	StudentStatus *StudentStatus `json:"student_status,omitempty"`
}

func (u *Users) ApplyEdit(e *EditUserRequest) {
	// TODO: Can be done better, please look into it.
	if e.Email != nil {
		u.Email = *e.Email
	}
	if e.Role != nil {
		u.Role = *e.Role
	}
	if e.FirstName != nil {
		u.FirstName = *e.FirstName
	}
	if e.LastName != nil {
		u.LastName = *e.LastName
	}
	if e.Gender != nil {
		u.Gender = *e.Gender
	}
	if e.PersonalNo != nil {
		u.PersonalNo = *e.PersonalNo
	}
	if e.PhoneNumber != nil {
		u.PhoneNumber = *e.PhoneNumber
	}
	if e.PhotoUrl != nil {
		u.PhotoUrl = *e.PhotoUrl
	}
	if e.StudentStatus != nil {
		statusCopy := *e.StudentStatus
		u.StudentStatus = &statusCopy
	}
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
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
