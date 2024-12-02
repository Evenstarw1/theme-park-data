package themepark

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type Handlers struct {
	db Store
}

func NewHandlers(db Store) *Handlers {
	return &Handlers{db: db}
}

func (h *Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello theme park world :)"))
}

func (h *Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"]) // Check error!!!
	user, err := h.db.GetUser(userID)
	if err != nil {
		w.Write([]byte("Todo mal..."))
	}

	// Marshall to JSON
	b, err := json.Marshal(user)
	if err != nil {
		w.Write([]byte("Todo mal..."))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *Handlers) SignIn(w http.ResponseWriter, r *http.Request) {
	var userLogin UserLogin
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, userId, err := h.db.SignIn(userLogin.Email, userLogin.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenResponse := Token{userId, token}
	tokenJSON, err := json.Marshal(tokenResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(tokenJSON)
}

func (h *Handlers) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]

		if h.db.IsLoggedIn(reqToken) {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

func (h *Handlers) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.db.GetAllCategories()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categoriesJson, err := json.Marshal(categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(categoriesJson)
}

func (h *Handlers) AddUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.AddUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	userID, _ := strconv.Atoi(vars["id"]) // Check error!!!
	user.ID = userID

	err = h.db.UpdateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
