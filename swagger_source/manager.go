package swagger_source

// @Summary Get student subject overview list
// @Description Retrieves a list of subjects, their groups and how many students are in them.
// @Tags manager
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.JsonGetSubjectStudentOverviewResponse
// @Failure 400 {object} models.JsonErrorResponse "Error returned will be 'Manager has no faculty attached.'"
// @Failure 500 {object} models.JsonErrorResponse "Error returned will be 'Failed to fetch.'"
// @Failure 403 {object} models.JsonErrorResponse "Error returned will be: 'Unauthorized from access.'"
// @Router /api/manager/subject-student-overview [get]
func _doc_GetSubjectStudentOverview() {}