package response

type DocsErrorResponse struct {
	Message string `json:"message" example:"Invalid request body"`
	Errors  struct {
		FieldName string `json:"fieldName" example:"error message"`
	} `json:"errors"`
}

type DocsErrorResponseInternalServerErr struct {
	Message string `json:"message" example:"internal server error"`
}

type DocsSuccessResponseLoginUser struct {
	Data struct {
		Token string `json:"token" example:"jwt.access.token"`
	} `json:"data"`
}

type DocsSuccessResponseRegisterUser struct {
	Message string `json:"message" example:"Done."`
}
