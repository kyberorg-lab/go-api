package user

type OldUser struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func GetSampleUser() OldUser {
	return OldUser{
		ID:       1,
		Username: "username",
		Password: "password",
	}
}
