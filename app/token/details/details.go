package details

type TokenDetails struct {
	AccessToken         string
	RefreshToken        string
	RefreshUuid         string
	UserAgent           string
	AccessTokenExpires  int64
	RefreshTokenExpires int64
	CreatedAt           int64
}
