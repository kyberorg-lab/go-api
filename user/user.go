package user

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetSampleUser() User {
	return User{
		ID:       1,
		Username: "username",
		Password: "password",
	}
}
