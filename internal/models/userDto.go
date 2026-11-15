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
	Email         string         `json:"email" validate:"required,email" example:"student@example.com"`
	Password      string         `json:"password" validate:"required,min=8" example:"securepassword1"`
	Role          Role           `json:"role" swaggertype:"string" example:"student"`
	FirstName     string         `json:"first_name" validate:"required" example:"Mariam"`
	LastName      string         `json:"last_name" validate:"required" example:"Mikava"`
	Gender        Gender         `json:"gender" swaggertype:"string" example:"active"`
	PersonalNo    string         `json:"personal_no" example:"01032010121"`
	PhoneNumber   string         `json:"phone_number" example:"+9955111092345"`
	PhotoUrl      string         `json:"photo_url" example:"http://link_to_photo"`
	StudentStatus *StudentStatus `json:"student_status,omitempty" swaggertype:"string" example:"active"`
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
	Email         *string        `json:"email,omitempty" example:"student@example.com"`
	Role          *Role          `json:"role,omitempty" swaggertype:"string" example:"student"`
	FirstName     *string        `json:"first_name,omitempty" example:"Mariam"`
	LastName      *string        `json:"last_name,omitempty" example:"Mikava"`
	Gender        *Gender        `json:"gender,omitempty" swaggertype:"string" example:"active"`
	PersonalNo    *string        `json:"personal_no,omitempty" example:"01032010121"`
	PhoneNumber   *string        `json:"phone_number,omitempty" example:"+9955111092345"`
	PhotoUrl      *string        `json:"photo_url,omitempty" example:"http://link_to_photo"`
	StudentStatus *StudentStatus `json:"student_status,omitempty" swaggertype:"string" example:"active"`
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
	Email    string `json:"email" validate:"required,email" example:"student@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"securepassword123"`
}

type LoginUserResponse struct {
	ID            uint           `json:"id" example:"1"`
	Email         string         `json:"email" example:"student@example.com"`
	Role          string         `json:"role" example:"1"`
	FirstName     string         `json:"first_name" example:"John"`
	LastName      string         `json:"last_name" example:"Doe"`
	PhotoUrl      string         `json:"photo_url" example:"http://link_to_image"`
	StudentStatus *StudentStatus `json:"student_status,omitempty" swaggertype:"string" example:"active"`
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
