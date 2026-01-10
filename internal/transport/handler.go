package transport

import (
	"log"
	"net/http"

	"github.com/AI-Hackathon-2026/Clients-Service/internal/model"
)

type service interface {
    findProfile(username string) (model.GetUserProfilePayload, error)
    updateProfile(username string) (model.UpdateUserProfilePayload, error)
    findGraph(graphId string) (model.GraphPayload, error)
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
    router.HandleFunc("PUT /profile-edit/{username}", h.putProfile)

    router.HandleFunc("GET /graphs/{username}", h.getGraphs)
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
    username := r.PathValue("username")
    prfl, err := h.srv.findProfile(username)
    if err != nil {
        http.Error(w, "Something went wrong", http.StatusBadRequest)
        log.Printf("getProfile error: %v", err)
        return
    }
    
}

func (h *Handler) putProfile(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getGraphs(w http.ResponseWriter, r *http.Request) {

}
