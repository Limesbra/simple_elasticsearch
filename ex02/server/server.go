package server

import (
	"encoding/json"
	"errors"
	elastic "ex02/db"
	"ex02/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	Store Store
}

type Store interface {
	GetPlaces(limit int, offset int) ([]model.Place, int, error)
}

func NewServer(store Store) *Server {
	return &Server{Store: store}
}

func (s *Server) handleApi(w http.ResponseWriter, r *http.Request) {
	pageSize := 10
	pageParam := r.URL.Query().Get("page")
	if pageParam == "" {
		pageParam = "1" // По умолчанию возвращаем первую страницу
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	offset := (page - 1) * pageSize

	a, total, err := s.Store.GetPlaces(pageSize, offset)

	if a == nil {
		err := errors.New("invalid 'page' value: 'foo'")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := total / pageSize

	if total%10 != 0 {
		totalPages++
	}

	if page <= 0 || page > totalPages {
		http.Error(w, "requested page is invalid", http.StatusBadRequest)
		log.Println("requested page is invalid", totalPages)
		return
	}

	// HTMLpage := index.BuildHTML(totalPages, pageSize, page, a)
	w.Header().Set("Content-Type", "application/json")

	meta := elastic.AddMetaInfo(a, total, page, totalPages)

	jsonData, err := json.MarshalIndent(meta, "", "  ")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(jsonData))

}

func (s *Server) Run() {
	http.HandleFunc("/api/places", s.handleApi)
	err := http.ListenAndServe(":8888", nil)

	if err != nil {
		log.Fatal(err)
	}
}
