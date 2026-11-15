package swagger_source

// @Summary Ping the server
// @Description Returns a Pong message to check if the server is alive
// @Tags system
// @Produce json
// @Success 200 {object} models.JsonSuccesfulResponse "Pong!"
// @Router /ping [get]
func _doc_Ping() {}

// @Summary Login
// @Description Authenticate a user and return a JWT token in an HttpOnly cookie
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginUserRequest true "Login Credentials"
// @Success 200 {object} models.JsonLoginResponse
// @Failure 400 {object} models.JsonErrorResponse "Error returned will be: 'Invalid request format'"
// @Failure 401 {object} models.JsonErrorResponse "Error returned will be: 'Invalid email or password'"
// @Router /api/login [post]
func _doc_Login() {}
