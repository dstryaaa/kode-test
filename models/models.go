package models

type User struct {
	ID       string `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Note struct {
	ID       string `json:"userID"`
	Username string `json:"username"`
	Note     string `json:"note"`
}
