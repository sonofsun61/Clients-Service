package transport

import (
	"log/slog"
	"net/http"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
	"github.com/AI-Hackathon-2026/Clients-Service/internal/service"
)

type ProfileService interface {
	FindProfile(username string) (*model.GetUserProfilePayload, error)
	UpdateProfile(prf *model.UpdateUserProfilePayload) error
	FindGraphs(username string) ([]model.GraphPayload, error)
}

type Handler struct {
	profileService ProfileService
	authService    service.AuthService
	logger         *slog.Logger
}

func NewHandler(profileService ProfileService, authService service.AuthService, logger *slog.Logger) *Handler {
	return &Handler{
		profileService: profileService,
		authService:    authService,
		logger:         logger,
	}
}

func (h *Handler) RegisterPublicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", h.postRegister)
	mux.HandleFunc("POST /auth/login", h.postLogin)

	mux.HandleFunc("POST /auth/refresh", h.postRefresh)
}

func (h *Handler) RegisterProtectedRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/logout", h.postLogout)

	mux.HandleFunc("GET /profile-edit/{username}", h.getProfile)
	mux.HandleFunc("PUT /profile-edit/", h.putProfile)

	mux.HandleFunc("GET /graphs/{username}", h.getGraphs)
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

func (h *Handler) postRegister(w http.ResponseWriter, r *http.Request) {
	req := new(model.RegisterRequest)
	if err := ParseJson(r, req); err != nil {
		raiseError(w, "postRegister, ParseJson error: ", err)
		return
	}
	if err := h.authService.Register(r.Context(), *req); err != nil {
		raiseError(w, "postRegister, authService.Register error: ", err)
		return
	}
	WriteJson(w, http.StatusCreated, req)
}

func (h *Handler) postLogin(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := ParseJson(r, &req); err != nil {
		raiseError(w, "invalid login payload", err)
		return
	}
	tokens, err := h.authService.Login(r.Context(), req)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	WriteJson(w, http.StatusOK, tokens)
}

func (h *Handler) postRefresh(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshRequest
	if err := ParseJson(r, &req); err != nil {
		raiseError(w, "invalid refresh payload", err)
		return
	}

	newTokens, err := h.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		raiseError(w, "refresh failed", err)
		return
	}

	WriteJson(w, http.StatusOK, newTokens)
}

func (h *Handler) postLogout(w http.ResponseWriter, r *http.Request) {
	var req model.RefreshRequest
	if err := ParseJson(r, &req); err != nil {
		raiseError(w, "invalid logout payload", err)
		return
	}
	err := h.authService.Logout(r.Context(), req.RefreshToken)
	if err != nil {
		raiseError(w, "logout failed", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
