package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/evenstarw1/theme-park-data/pkg/themepark"
)

const (
	envDbConnectionString = "https://amsmpjefzegchbqafouw.supabase.co"
	envListeningPort      = "THEMEPARK_LISTENING_PORT"
)

type Config struct {
	databaseConnectionString string
	listeningPort            int
}

func main() {
	config := getEnv()

	router := mux.NewRouter()
	databaseStore, err := themepark.NewDatabaseStore(config.databaseConnectionString)
	if err != nil {
		os.Exit(1)
	}

	defer databaseStore.Close()

	handlers := themepark.NewHandlers(databaseStore)

	pub := router.PathPrefix("/pub").Subrouter()
	pub.HandleFunc("/login", handlers.SignIn).
		Methods(http.MethodPost)
	pub.HandleFunc("/getCategories", handlers.GetCategories).
		Methods(http.MethodGet)
	pub.HandleFunc("/register", handlers.AddUser).
		Methods(http.MethodPost)

	priv := router.PathPrefix("/priv").Subrouter()
	priv.Use(handlers.AuthMiddleware)
	priv.HandleFunc("/", handlers.ServeHTTP).
		Methods(http.MethodGet)
	priv.HandleFunc("/users/{id:[0-9]+}", handlers.GetUser).
		Methods(http.MethodGet)
	priv.HandleFunc("/users/{id:[0-9]+}", handlers.UpdateUser).
		Methods(http.MethodPatch)
	priv.HandleFunc("/parks", handlers.GetParks).
		Methods(http.MethodGet)
	priv.HandleFunc("/parks", handlers.InsertThemePark).
		Methods(http.MethodPost)
	priv.HandleFunc("/park/{id:[0-9]+}", handlers.GetParkDetails).
		Methods(http.MethodGet)
	priv.HandleFunc("/park/{id:[0-9]+}", handlers.DeleteThemePark).
		Methods(http.MethodDelete)
	priv.HandleFunc("/park/{id:[0-9]+}", handlers.UpdateThemePark).
		Methods(http.MethodPatch)
	priv.HandleFunc("/park/comments", handlers.InsertParkComment).
		Methods(http.MethodPost)
	priv.HandleFunc("/categories", handlers.AddCategory).
		Methods(http.MethodPost)
	priv.HandleFunc("/users", handlers.GetAllUsers).
		Methods(http.MethodGet)

	http.ListenAndServe(fmt.Sprintf(":%d", config.listeningPort), router)
}

func getEnv() Config {
	var config Config
	config.databaseConnectionString = os.Getenv(envDbConnectionString)
	if config.databaseConnectionString == "" {
		config.databaseConnectionString = "postgres://localhost:5432/themepark"
	}

	portStr := os.Getenv(envListeningPort)
	if portStr == "" {
		config.listeningPort = 8080
	} else {
		config.listeningPort, _ = strconv.Atoi(portStr)
	}

	return config
}
