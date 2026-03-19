package transport

import (
	"net/http"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type ProfileService interface {
	FindProfile(username string) (*model.GetUserProfilePayload, error)
	UpdateProfile(prf *model.UpdateUserProfilePayload) error
	FindGraphs(username string) ([]model.GraphPayload, error)
}

type AuthService interface {
	Register(req model.RegisterRequest) error
	Login(req model.LoginRequest) (*model.AuthResponse, error)
	Refresh(req model.RefreshRequest) (*model.AuthResponse, error)
	Logout(refreshToken string) error
}

type Handler struct {
	profileService ProfileService
	authService    AuthService
}

func NewHandler(profileService ProfileService, authService AuthService) *Handler {
	return &Handler{
		profileService: profileService,
		authService:    authService,
	}
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /profile-edit/{username}", h.getProfile)
	router.HandleFunc("PUT /profile-edit/", h.putProfile)

	router.HandleFunc("GET /graphs/{username}", h.getGraphs)

	router.HandleFunc("POST /auth/register", h.postRegister)
	router.HandleFunc("POST /auth/login", h.postLogin)
	router.HandleFunc("POST /auth/refresh", h.postRefresh)
	router.HandleFunc("POST /auth/logout", h.postLogout)
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	prfl, err := h.profileService.FindProfile(username)
	if err != nil {
		raiseError(w, "getProfile error:", err)
		return
	}
	WriteJson(w, http.StatusOK, prfl)
}

func (h *Handler) putProfile(w http.ResponseWriter, r *http.Request) {
	prfl := new(model.UpdateUserProfilePayload)
	if err := ParseJson(r, prfl); err != nil {
		raiseError(w, "putProfile, ParseJson error:", err)
		return
	}
	if err := h.profileService.UpdateProfile(prfl); err != nil {
		raiseError(w, "putProfile, updateProfile error:", err)
		return
	}
}

func (h *Handler) getGraphs(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("username")
	graphs, err := h.profileService.FindGraphs(userId)

	if err != nil {
		raiseError(w, "getGraphs, FindGraphs error:", err)
		return
	}
	WriteJson(w, http.StatusOK, graphs)
}

// Авторизация
func (h *Handler) postRegister(w http.ResponseWriter, r *http.Request) {
	req := new(model.RegisterRequest)
	if err := ParseJson(r, req); err != nil {
		raiseError(w, "postRegister, ParseJson error: ", err)
		return
	}
	if err := h.authService.Register(*req); err != nil {
		raiseError(w, "postRegister, authService.Register error: ", err)
		return
	}
	WriteJson(w, http.StatusCreated, req)
}

func (h *Handler) postLogin(w http.ResponseWriter, r *http.Request) {
	req := new(model.LoginRequest)
	if err := ParseJson(r, req); err != nil {
		raiseError(w, "postLogin, ParseJson error:", err)
		return
	}
	tokens, err := h.authService.Login(*req)
	if err != nil {
		raiseError(w, "postLogin, authService.Login error:", err)
		return
	}
	WriteJson(w, http.StatusOK, tokens)
}

func (h *Handler) postRefresh(w http.ResponseWriter, r *http.Request) {
	req := new(model.RefreshRequest)
	if err := ParseJson(r, req); err != nil {
		raiseError(w, "postRefresh, ParseJson error:", err)
		return
	}
	tokens, err := h.authService.Refresh(*req)
	if err != nil {
		raiseError(w, "postRefresh, authService.Refresh error:", err)
		return
	}

	WriteJson(w, http.StatusOK, tokens)
}

func (h *Handler) postLogout(w http.ResponseWriter, r *http.Request) {
	req := new(model.RefreshRequest)
	if err := ParseJson(r, req); err != nil {
		raiseError(w, "postLogout, ParseJson error:", err)
		return
	}
	if err := h.authService.Logout(req.RefreshToken); err != nil {
		raiseError(w, "postLogout, authService.Logout error:", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
