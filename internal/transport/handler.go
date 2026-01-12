package transport

import (
	"net/http"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type service interface {
    FindProfile(username string) (*model.GetUserProfilePayload, error)
    UpdateProfile(prf *model.UpdateUserProfilePayload) error
    FindGraphs(username string) ([]model.GraphPayload, error)
}

type Handler struct {
    srv service
}

func NewHandler(srv service) *Handler {
    return &Handler{
        srv: srv,
    }
}

func (h *Handler) RegisterRoutes(router *http.ServeMux) {
    router.HandleFunc("GET /profile-edit/{username}", h.getProfile)
    router.HandleFunc("PUT /profile-edit/", h.putProfile)

    router.HandleFunc("GET /graphs/{username}", h.getGraphs)
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
    username := r.PathValue("username")
    prfl, err := h.srv.FindProfile(username)
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
    if err := h.srv.UpdateProfile(prfl); err != nil {
        raiseError(w, "putProfile, updateProfile error:", err)
        return
    }
}

func (h *Handler) getGraphs(w http.ResponseWriter, r *http.Request) {
    userId := r.PathValue("username")    
    graphs, err := h.srv.FindGraphs(userId)

    if err != nil {
        raiseError(w, "getGraphs, FindGraphs error:", err)
        return
    }
    WriteJson(w, http.StatusOK, graphs)
}

