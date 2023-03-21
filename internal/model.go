package togo

type UserCredentials struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type UserSignupRequest struct {
	UserCredentials
}

type UserLoginRequest struct {
	UserCredentials
}
