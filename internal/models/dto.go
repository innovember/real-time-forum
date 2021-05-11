package models

type InputUserSignIn struct {
	Nickname string `json:"usernameOrEmail"` // username or email
	Password string `json:"password"`
}

type InputUserSignUp struct {
	Nickname  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
}

type InputPost struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	Categories []string `json:"categories"`
}
