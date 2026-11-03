package models

type LoginUserDTO struct {
    ID            uint          	`json:"id"`
    Email         string        	`json:"email"`
    Role          Role          	`json:"role"`
    FirstName     string        	`json:"first_name"`
    LastName      string        	`json:"last_name"`
    PhotoUrl      string        	`json:"photo_url"`
    StudentStatus *StudentStatus 	`json:"student_status,omitempty"`
}