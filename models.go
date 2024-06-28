package main

type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
}

type Vote struct {
    UserID string `json:"user_id"`
    Choice string `json:"choice"`
}
