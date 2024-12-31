package themepark

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

const (
	AdminAccessLevel = 1
	UserAccessLevel  = 2
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

func (h *Handlers) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	userAccess, _ := strconv.Atoi(r.Header.Get("app-user-access-level")) // Check error!!!

	if userAccess == UserAccessLevel {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	users, err := h.db.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	b, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

		isLoggedIn, userId, accessLevel := h.db.IsLoggedIn(reqToken)
		if isLoggedIn {
			r.Header.Set("app-user-id", strconv.Itoa(userId))
			r.Header.Set("app-user-access-level", strconv.Itoa(accessLevel))
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

func (h *Handlers) AddCategory(w http.ResponseWriter, r *http.Request) {
	userAccess, _ := strconv.Atoi(r.Header.Get("app-user-access-level")) // Check error!!!

	if userAccess == UserAccessLevel {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var category Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.AddCategory(category.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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

func (h *Handlers) GetParks(w http.ResponseWriter, r *http.Request) {
	themeparks, err := h.db.GetAllThemeParks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	themeparksJson, err := json.Marshal(themeparks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(themeparksJson)
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GetParkDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	themeparkId, _ := strconv.Atoi(vars["id"]) // Check error!!!

	themepark, err := h.db.GetThemeParkDetail(themeparkId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	themeparkDetailJson, err := json.Marshal(themepark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(themeparkDetailJson)
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) InsertParkComment(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userId, _ := strconv.Atoi(r.Header.Get("app-user-id")) // Check error!!!

	// Get comment from json
	var comment Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.InsertParkComment(comment.ThemeparkId, userId, comment.Comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) InsertThemePark(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userAccess, _ := strconv.Atoi(r.Header.Get("app-user-access-level")) // Check error!!!

	if userAccess == UserAccessLevel {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get themePark from json
	var themePark ThemePark
	err := json.NewDecoder(r.Body).Decode(&themePark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert theme park
	err = h.db.AddThemePark(themePark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) DeleteThemePark(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userAccess, _ := strconv.Atoi(r.Header.Get("app-user-access-level")) // Check error!!!

	if userAccess == UserAccessLevel {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	themeParkId, _ := strconv.Atoi(vars["id"]) // Check error!!!

	// Delete theme park
	err := h.db.DeleteThemePark(themeParkId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handlers) UpdateThemePark(w http.ResponseWriter, r *http.Request) {
	// Get user id
	userAccess, _ := strconv.Atoi(r.Header.Get("app-user-access-level")) // Check error!!!

	if userAccess == UserAccessLevel {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	themeParkId, _ := strconv.Atoi(vars["id"]) // Check error!!!

	// Get themePark from json
	var themePark ThemePark
	err := json.NewDecoder(r.Body).Decode(&themePark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set id to themepark struct
	themePark.Id = themeParkId

	err = h.db.UpdateThemePark(themePark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
