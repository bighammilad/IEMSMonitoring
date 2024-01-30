package model

type LoginResult struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Role     int    `json:"role,omitempty"`
}

type UserModel struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Role     int    `json:"role,omitempty"`
}

type UserRes struct {
	UserId   int    `json:"userid,omitempty"`
	Username string `json:"username,omitempty"`
	Role     int    `json:"role,omitempty"`
}

type UserRegister struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Role     int    `json:"role,omitempty"`
}

type UserUpdate struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Role     int    `json:"role,omitempty"`
}

const (
	Admin AccessLevel = iota + 1
	Regular
	Demo
)

type AccessLevel int

type User struct {
	Username    string      `json:"username,omitempty"`
	AccessLevel AccessLevel `json:"access_level,omitempty"`
}

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}
