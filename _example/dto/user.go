package dto

// UserSignUp авторизация нового пользователя
type UserSignUp struct {
	FirstName       string `json:"first_name,omitempty" label:"ru:Имя;en:Name;fr:Nom"`
	LastName        string `json:"last_name,omitempty"  label:"ru:Фамилия;en:Surname;fr:Nom de famille"`
	Age             int    `json:"age,omitempty" validate:"example" label:"ru:Возраст;en:Age;fr:Âge"`
	Email           string `json:"email,omitempty" validate:"required,email" label:"ru:Почта;en:Email;fr:E-mail"`
	Password        string `json:"password,omitempty" validate:"required" label:"ru:Пароль;en:Password;fr:Mot de passe"`
	PasswordConfirm string `json:"password_confirm,omitempty" validate:"required,eqfield=Password" label:"ru:Подтверждение пароля;en:Password Confirm;fr:Confirmer le mot de passe"`
}
