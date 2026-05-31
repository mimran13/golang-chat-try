package handler

import (
	"database/sql"
	"go-chat/internal/response"
	"net/http"
)

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	if err := h.db.PingContext(r.Context()); err != nil {
		response.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{                                                                                                                       
              "status": "unhealthy",                                                                                                                                                                    
              "error":  "database unreachable",
          })                                                                                                                                                                                              
          return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{                                                                                                                                             
          "status": "healthy",                                                                                                                                                                          
      })    
}
