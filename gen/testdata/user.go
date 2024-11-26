package testdata

import "database/sql"

type User struct {
	Name     string
	Age      *int
	NickName *sql.NullString
	Picture  []byte
}

type UserDetail struct {
	Address string
}
