package togo

type omit *struct{}

type User struct {
	Username       *string `db:"username"`
	ID             *int64  `db:"id"`
	HashedPassword *string `db:"hashed_password"`
	PasswordSalt   *string `db:"password_salt"`
	CreatedAt      *string `db:"created_at"`
	UpdatedAt      *string `db:"updated_at"`
}

type UserCredentials struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

// The reason we do this is that we don't have to convert the User struct to
// UserInfo struct when we want to return the user info to the client.
// You may ask why we don't just remove the composite User struct and use
// only UserInfo struct instead. The reason is that when destructuring
// the user info from the database, the driver requires that the struct
// fields must match ALL the column names in the database.
// So the way we approach this is to have a composite struct, and if you want
// to reveal or obfusticate some fields to the client, you can just add the
// json tag to it. You'll see this pattern in other structs as well.
type UserInfo struct {
	*User

	Username       *string `json:"username" db:"username"`
	ID             *int64  `json:"-" db:"id"`
	HashedPassword *string `json:"-" db:"hashed_password"`
	PasswordSalt   *string `json:"-" db:"password_salt"`
	CreatedAt      *string `json:"created_at" db:"created_at"`
	UpdatedAt      *string `json:"-" db:"updated_at"`

	HasSuperPower *bool `json:"has_super_power,omitempty"`
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
