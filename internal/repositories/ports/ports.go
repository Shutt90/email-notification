package ports

type UserRepo interface {
	CreateUser()
	AuthenticateUser()
}
