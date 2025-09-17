package requests

type RegisterUserRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
}

type LoginUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=8"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"eqfield=NewPassword"`
}

type ForgotPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	FrontendURL string `json:"frontendUrl" validate:"required,url"`
}

type ResetPasswordRequest struct {
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"eqfield=Password"`
	Token           string `json:"token" validate:"required,min=5,max=6"`
	Meta            string `json:"meta" validate:"required"`
}
