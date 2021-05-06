package models

type InputUserSignIn struct {
	Nickname string `json:"nickname"` // username or email
	Password string `json:"password"`
}

type InputUserSignUp struct {
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
}
