package model

type GetUserProfilePayload struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateUserProfilePayload struct {
	UsernameNew string `json:"username_new"`
	UsernameOld string `json:"username_old"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type GraphPayload struct {
	GraphId string `json:"graph_id"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
