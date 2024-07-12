package server

import (
	"encoding/json"
	elastic "ex04/db"
	"ex04/model"
	"ex04/token"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	Store Store
}

type Store interface {
	GetPlaces(limit, offset float64) ([]model.Place, error)
}

func NewServer(store Store) *Server {
	return &Server{Store: store}
}

func (s *Server) handleApi(w http.ResponseWriter, r *http.Request) {

	tkn, err := token.ExtractTokenFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = token.ValidateToken(tkn)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	firstParam := r.URL.Query().Get("lat")
	secondParam := r.URL.Query().Get("lon")
	if firstParam == "" || secondParam == "" {
		http.Error(w, "Invalid parameters", http.StatusBadRequest)
		return
	}

	lat, errLat := strconv.ParseFloat(firstParam, 64)
	lon, errLon := strconv.ParseFloat(secondParam, 64)
	if errLat != nil || errLon != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	a, err := s.Store.GetPlaces(lat, lon)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	meta := elastic.AddInfo(a)

	jsonData, err := json.MarshalIndent(meta, "", "  ")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(jsonData))

}

func (s *Server) Run() {
	http.HandleFunc("/api/get_token", token.GetToken)
	http.HandleFunc("/api/recommend", s.handleApi)
	err := http.ListenAndServe(":8888", nil)

	if err != nil {
		log.Fatal(err)
	}
}
