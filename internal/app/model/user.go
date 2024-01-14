package model

const UserTableName = "users"

type User struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
	Age  string `db:"age"`
}
