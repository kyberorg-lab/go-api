package details

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	UserAgent    string
	AtExpires    int64
	RtExpires    int64
	CreatedAt    int64
}
