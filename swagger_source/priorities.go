package swagger_source

// @Summary Get priorities
// @Description Retrieves a list of priorities for the logged-in student
// @Tags priority
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.JsonGetPrioritiesResponse
// @Failure 401 {object} models.JsonErrorResponse "Error returned will be: 'Unauthorized from access.'"
// @Failure 500 {object} models.JsonErrorResponse
// @Router /api/priority [get]
func _doc_GetPriorities() {}

// @Summary Add priority
// @Description Creates a new priority for the logged-in student
// @Tags priority
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param priority body models.PriorityRequest true "Priority Data"
// @Success 201 {object} models.JsonSuccesfulResponse
// @Failure 400 {object} models.JsonErrorResponse
// @Failure 401 {object} models.JsonErrorResponse "Error returned will be: 'Unauthorized.'"
// @Router /api/priority [post]
func _doc_AddPriority() {}

// @Summary Delete priority
// @Description Deletes a priority by ID and re-indexes the remaining priorities
// @Tags priority
// @Produce json
// @Security BearerAuth
// @Param id path string true "Priority ID"
// @Success 200 {object} models.JsonSuccesfulResponse "Priority deleted and re-indexed"
// @Failure 400 {object} models.JsonErrorResponse "Service error"
// @Router /api/priority/{id} [delete]
func _doc_DeletePriority() {}
