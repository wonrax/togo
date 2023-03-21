package togo

type omit *struct{}

type UserCredentials struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

type Todo struct {
	ID          *int64  `json:"id" db:"id"`
	Owner       *int64  `json:"owner" db:"owner"`
	Title       *string `json:"title,omitempty" db:"title"`
	Description *string `json:"description,omitempty" db:"description"`
	Completed   *bool   `json:"completed,omitempty" db:"completed"`
	CreatedAt   *string `json:"created_at" db:"created_at"`
	UpdatedAt   *string `json:"updated_at" db:"updated_at"`
}

type UserSignupRequest struct {
	*UserCredentials
}

type UserLoginRequest struct {
	*UserCredentials
}

type TodoRequest struct {
	*Todo

	ID omit `json:"id,omitempty"`
}

type TodoResponse struct {
	*Todo
}

type ListTodosResponse struct {
	Todos []TodoResponse `json:"todos"`
}
