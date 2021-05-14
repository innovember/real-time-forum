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

type InputComment struct {
	PostID  int64  `json:"postID"`
	Content string `json:"content"`
}

type InputGetComments struct {
	Option        string `json:"option"` // user or post
	PostID        int64  `json:"post_id"`
	UserID        int64  `json:"user_id"`
	Offset        int    `json:"offset"`
	Limit         int    `json:"limit"`
	LastCommentID int    `json:"lastCommentID"`
}

type InputGetPosts struct {
	Option     string   `json:"option"` // all, categories or author
	AuthorID   int64    `json:"authorID"`
	Categories []string `json:"categories"`
	Offset     int      `json:"offset"`
	Limit      int      `json:"limit"`
	LastPostID int      `json:"lastPostID"`
}
