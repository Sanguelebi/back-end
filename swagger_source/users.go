package swagger_source

// @Summary Get all users
// @Description Retrieves a list of all users. Requires Admin privileges.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.JsonGetUsersResponse
// @Failure 500 {object} models.JsonErrorResponse
// @Router /api/users [get]
func _doc_GetAllUsers() {}

// @Summary Add a new user
// @Description Creates a new user in the system. Requires Admin privileges.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body models.RegisterUserRequest true "User Data"
// @Success 201 {object} models.JsonAddUserResponse
// @Failure 400 {object} models.JsonErrorResponse
// @Router /api/user [post]
func _doc_AddUser() {}

// @Summary Edit a user
// @Description Updates an existing user by ID. Requires Admin privileges.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param user body models.EditUserRequest true "Updated User Data"
// @Success 200 {object} models.JsonSuccesfulResponse "User edited succesfully!"
// @Failure 400 {object} models.JsonErrorResponse
// @Router /api/user/{id} [patch]
func _doc_EditUser() {}

// @Summary Delete a user
// @Description Deletes a user by ID. Requires Admin privileges.
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} models.JsonSuccesfulResponse "User deleted succesfully!"
// @Failure 400 {object} models.JsonErrorResponse
// @Router /api/user/{id} [delete]
func _doc_DeleteUser() {}
