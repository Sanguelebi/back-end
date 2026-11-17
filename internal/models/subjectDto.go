package models

type PriorityRequest struct {
	SubjectGroupID uint               `json:"subject_group_id" example:"2"`
	SlotID         int                `json:"slot_id" example:"1"`
	SubjectType    SubjectType 		  `json:"subject_type" swaggertype:"string" example:"free"`
}

type ChoosableSubjectResponse struct {
	SubjectData   Subject        `json:"subjectData"`
	SubjectGroups []SubjectGroup `json:"subjectGroups"`
}

type SubjectGroupOverview struct {
	GroupID          uint   `json:"group_id"`
	GroupNo          int16  `json:"group_no"`
	StudentIDs       []uint `json:"student_ids"`
	AmountOfStudents int    `json:"amount_of_students"`
	MaxAmount        int8   `json:"max_amount"`
}

type SubjectStudentOverviewResponse struct {
	Subject       Subject         `json:"subject"`
	SubjectGroups []SubjectGroupOverview `json:"subject_groups"`
}