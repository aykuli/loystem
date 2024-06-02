package repository

type Repository interface {
	CreateUsersTable() error
}
