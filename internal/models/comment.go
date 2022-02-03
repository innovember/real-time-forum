package models

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"postId"`
	AuthorID  int64  `json:"-"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"createdAt"`
	Author    *User  `json:"author"`
}
