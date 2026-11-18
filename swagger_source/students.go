package swagger_source

// @Summary Get choosable subjects for a student
// @Description Retrieves a list of subjects that a student can choose
// @Tags student
// @Produce json
// @Security BearerAuth
// @Param subjectType query string true "Filter by subject type" Enums(professional, free)
// @Success 200 {array} models.JsonChoosableSubjectResponse
// @Failure 400 {object} models.JsonErrorResponse "Error returned will be 'Specify subject type.'"
// @Failure 403 {object} models.JsonErrorResponse "Error returned will be: 'Unauthorized from access.'"
// @Router /api/student/choosable-subject [get]
func _doc_GetStudentChoosableSubjects() {}

// @Summary Get student subject allowance
// @Description Retrieves the subject slot allowances assigned to the authenticated student.
// @Tags student
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.JsonSubjectAllowanceResponse
// @Failure 401 {object} models.JsonErrorResponse "Error returned will be: 'Unauthorized from access.'"
// @Failure 500 {object} models.JsonErrorResponse "Error returned will be 'failed to fetch student subject allowance'"
// @Router /api/student/subject-allowance [get]
func _doc_GetStudentSubjectAllowance() {}

// @Summary Get priorities
// @Description Retrieves a list of priorities for the logged-in student
// @Tags student
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.JsonGetPrioritiesResponse
// @Failure 403 {object} models.JsonErrorResponse "Error returned will be: 'Unauthorized from access.'"
// @Failure 500 {object} models.JsonErrorResponse
// @Router /api/priority [get]
func _doc_GetPriorities() {}

// @Summary Add priority
// @Description Creates a new priority for the logged-in student
// @Tags student
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param priority body models.PriorityRequest true "Priority Data"
// @Success 201 {object} models.JsonSuccesfulResponse
// @Failure 400 {object} models.JsonErrorResponse
// @Failure 403 {object} models.JsonErrorResponse "Error returned will be: 'Unauthorized.'"
// @Router /api/priority [post]
func _doc_AddPriority() {}

// @Summary Delete priority
// @Description Deletes a priority by ID and re-indexes the remaining priorities
// @Tags student
// @Produce json
// @Security BearerAuth
// @Param id path string true "Priority ID"
// @Success 200 {object} models.JsonSuccesfulResponse "Priority deleted and re-indexed"
// @Failure 400 {object} models.JsonErrorResponse "Service error"
// @Router /api/priority/{id} [delete]
func _doc_DeletePriority() {}
