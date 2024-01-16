package dto

type User struct {
	Name string `json:"name" validate:"required" label:"ru:Имя;en:Name;fr:Nom"`
	Age  int    `json:"age" validate:"required,example" label:"ru:Возраст;en:Age;fr:Ajy"`
}
