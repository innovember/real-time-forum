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
	PostID  int64  `json:"postId"`
	Content string `json:"content"`
}

type InputGetComments struct {
	Option        string `json:"option"` // user or post
	PostID        int64  `json:"postId"`
	UserID        int64  `json:"userId"`
	Offset        int    `json:"offset"`
	Limit         int    `json:"limit"`
	LastCommentID int    `json:"lastCommentID"`
}

type InputGetPosts struct {
	Option     string   `json:"option"` // all, categories or author
	AuthorID   int64    `json:"authorId"`
	Categories []string `json:"categories"`
	Offset     int      `json:"offset"`
	Limit      int      `json:"limit"`
	LastPostID int      `json:"lastPostID"`
}

type InputRoom struct {
	RoomID        int64 `json:"roomId"`
	UserID        int64 `json:"userId"`
	MessageID     int64 `json:"messageId"`
	LastMessageID int64 `json:"lastMessageID"`
}
