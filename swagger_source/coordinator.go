package swagger_source

// @Summary Get coordinator student priorities
// @Description Retrieves a sorted list of all student priorities belonging to the coordinator's assigned faculty. Data is naturally grouped by student ID, slot ID, and priority.
// @Tags coordinator
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.JsonGetCoordinatorPrioritiesResponse "Returns a message and a 'data' array containing models.StudentPriorityList"
// @Failure 400 {object} models.JsonErrorResponse "Error returned will be 'coordinator is not assigned to any faculty'"
// @Failure 401 {object} models.JsonErrorResponse "Error returned will be 'Unauthorized from access.'"
// @Failure 404 {object} models.JsonErrorResponse "Error returned will be 'coordinator with id X not found'"
// @Failure 500 {object} models.JsonErrorResponse "Error returned will be 'failed to fetch student priorities' or other server errors."
// @Router /api/coordinator/student-priorities [get]
func _doc_GetCoordinatorStudentPriorities() {}


// @Summary Confirm or unconfirm a student's subject priority
// @Description Toggles the confirmation status of a specific student priority. If confirming, it automatically removes any other confirmed priority in the same slot. Checks group capacity before confirming.
// @Tags coordinator
// @Produce json
// @Security BearerAuth
// @Param id query int true "The ID of the Student Priority List record to confirm/unconfirm"
// @Success 200 {object} models.JsonSuccesfulResponse "Message returned will be 'Coordinator changed confirm succesfully!'"
// @Failure 400 {object} models.JsonErrorResponse "Error returned will be 'No priority ID Specified.' or parsing error."
// @Failure 401 {object} models.JsonErrorResponse "Error returned will be 'Unauthorized from access.'"
// @Failure 404 {object} models.JsonErrorResponse "Error returned will be 'priority record with id X not found'"
// @Failure 409 {object} models.JsonErrorResponse "Error returned will be 'subject group X is full (Y/Z)'"
// @Failure 500 {object} models.JsonErrorResponse "Error returned will be related to database transactions, fetching, or locks."
// @Router /api/coordinator/confirm-subject [post]
func _doc_CoordinatorConfirmSubject() {}