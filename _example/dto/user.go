package dto

type UserSignUp struct {
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	Email           string `json:"email,omitempty" validate:"required,email" name:"[ru:Электронная почта;en:Email;fr:Adresse e-mail]"`
	Password        string `json:"password,omitempty" validate:"required" name:"[ru:Пароль;en:Password;fr:Mot de passe]"`
	PasswordConfirm string `json:"password_confirm,omitempty" validate:"required,eqfield=Password" name:"[ru:Подтверждение пароля;en:Password confirmation;fr:Confirmation mot de passe]"`
}

type UserSignIn struct {
	Email    string `json:"email,omitempty" validate:"required,email" name:"[ru:Электронная почта;en:Email;fr:Adresse e-mail]"`
	Password string `json:"password,omitempty" validate:"required" name:"[ru:Пароль;en:Password;fr:Mot de passe]"`
}
