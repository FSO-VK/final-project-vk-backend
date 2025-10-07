package application

type AuthApplication struct {
	CheckAuth    CheckAuth
	LoginByEmail LoginByEmail
	Logout       Logout
}
