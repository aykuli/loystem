package user

type User struct {
	ID             string `json:"id"`
	Login          string `json:"login"`
	HashedPassword string `json:"-"`
	Salt           string
}
