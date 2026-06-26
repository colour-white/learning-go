package api

import (
	"database/sql"
	"encoding/json"
	"go-status/internal/models"
	"go-status/internal/monitor"
	"net/http"
)

type Server struct {
	DB *sql.DB
	Manager *monitor.Manager
}

// GET /targets
func (s *Server) listTargets(w http.ResponseWriter, r*http.Request){
	targets,err:=models.SelectAllTargets(s.DB)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJSON(w, http.StatusOK, targets)
}


// GET /probes
func (s *Server) listProbes(w http.ResponseWriter, r*http.Request){
	targets,err:=models.SelectAllProbes(s.DB)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeJSON(w, http.StatusOK, targets)
}



func writeJSON(w http.ResponseWriter, status int, v any){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (s *Server) Routes() *http.ServeMux {
      mux := http.NewServeMux()
      mux.HandleFunc("GET /targets", s.listTargets)
      mux.HandleFunc("GET /probes", s.listProbes)
      return mux
}
