package model

const UserTableName = "users"

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
}
