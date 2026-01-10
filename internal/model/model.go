package model

type GetUserProfilePayload struct {
    Id string `json:"id"`
    Username string `json:"username"`
    Email string `json:"email"` 
}

type UpdateUserProfilePayload struct {
    Username string `json:"username"`
    Email string `json:"email"`
    Password string `json:"password"` 
}

type GraphPayload struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Preview string `json:"preview"`
}
