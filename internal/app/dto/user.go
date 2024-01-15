package dto

type User struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required"`
}
