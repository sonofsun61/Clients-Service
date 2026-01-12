package model

type GetUserProfilePayload struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateUserProfilePayload struct {
	UsernameNew string `json:"username_new"`
    UsernameOld string `json:"username_old"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GraphPayload struct {
    GraphId int64 `json:"graph_id"`
}
