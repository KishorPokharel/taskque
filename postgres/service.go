package postgres

import "database/sql"

type Service struct {
	User  UserService
	Token TokenService
}

func NewService(db *sql.DB) Service {
	s := Service{
		User:  UserService{DB: db},
		Token: TokenService{DB: db},
	}
	return s
}
