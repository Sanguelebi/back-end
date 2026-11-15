package models

type JsonErrorResponse struct {
    Error string `json:"error" example:"Invalid request format"`
}

type JsonSuccesfulResponse struct {
    Message string `json:"message" example:"The action has been executed succesfully!"`
}

// ==== USER =====

type JsonGetUsersResponse struct {
    Message string `json:"message" example:"Fetched users succesfully!"`
    Users   []Users `json:"users"`
}

type JsonAddUserResponse struct {
    Message string `json:"message" example:"User added succesfully!"`
    ID      uint   `json:"id" example:"2"`
}

type JsonLoginResponse struct {
    Message string `json:"message" example:"User logged in succesfully!"`
    User LoginUserResponse `json:"user"`
}

// ==== SUBJECTS =====

type JsonGetSubjectsResponse struct {
    Message string `json:"message" example:"Subjects fetched succesfully!"`
    Subjects []Subject `json:"subjects"`
}

// ==== PRIORITIES ====

type JsonGetPrioritiesResponse struct {
    Message string `json:"message" example:"Priority list fetched succesfully!"`
    Priorities []StudentPriorityList `json:"priorities"`
}