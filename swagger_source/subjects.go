package swagger_source

// @Summary Get all subjects
// @Description Retrieves a list of all available subjects, with optional query filters
// @Tags subjects
// @Produce json
// @Param faculty query string false "Filter by faculty"
// @Param name query string false "Filter by subject name"
// @Success 200 {array} models.JsonGetSubjectsResponse
// @Failure 500 {object} models.JsonErrorResponse
// @Router /api/subjects [get]
func _doc_GetSubjects() {}
