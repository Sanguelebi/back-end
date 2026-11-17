package models

type PriorityRequest struct {
	SubjectGroupID uint  `json:"subject_group_id" example:"2"`
}

type ChoosableSubjectResponse struct {
	SubjectData   Subject        `json:"subjectData"`
	SubjectGroups []SubjectGroup `json:"subjectGroups"`
}