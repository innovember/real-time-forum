package models

type User struct {
	ID         int64  `json:"id,omitempty"`
	Nickname   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"-"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Age        int    `json:"age,omitempty"`
	Gender     string `json:"gender,omitempty"`
	CreatedAt  int64  `json:"createdAt,omitempty"`
	LastActive int64  `json:"lastActive,omitempty"`
	Status     string `json:"status"`
}
