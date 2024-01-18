package dto

// UserSignup авторизация нового пользователя
type UserSignup struct {
	FirstName       string `json:"first_name,omitempty" name:"Имя"`
	LastName        string `json:"last_name,omitempty"  name:"Фамилия"`
	Password        string `json:"password,omitempty" validate:"required" name:"[ru:Пароль;en:Password;fr:Mot de passe]"`
	PasswordConfirm string `json:"password_confirm,omitempty" validate:"required,eqfield=Password" name:"[ru:Подтверждение пароля;en:Password confirmation;fr:Confirmation mot de passe]"`
}
