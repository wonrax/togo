package togo

type UserSignupRequest struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}
