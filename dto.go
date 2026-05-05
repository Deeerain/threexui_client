package threexuiclinet

type LoginRequest struct {
	Username      string  `json:"username"`
	Password      string  `json:"password"`
	TwoFactorCode *string `json:"twoFactorCode"`
}

func NewLoginRequest(username string, password string, twoFactorCode *string) *LoginRequest {
	return &LoginRequest{
		Username:      username,
		Password:      password,
		TwoFactorCode: twoFactorCode,
	}
}
