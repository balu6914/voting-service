package main

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Vote represents a vote cast by a user
type Vote struct {
	UserID string `json:"user_id"`
	Choice string `json:"choice"`
}
